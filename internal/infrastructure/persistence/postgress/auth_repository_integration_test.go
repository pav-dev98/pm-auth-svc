package postgress_test

// Los tests de integración verifican que nuestro código funciona correctamente
// contra una base de datos real. A diferencia de los tests unitarios (que usan
// mocks/fakes), aquí levantamos un PostgreSQL de verdad usando Docker.
//
// Herramienta: testcontainers-go
// Levanta un contenedor Docker efímero durante la ejecución del test.
// El contenedor arranca, corremos los tests, y se destruye al finalizar.
// Esto garantiza aislamiento total: cada ejecución parte de cero.

import (
	"context"
	"testing"
	"time"

	// El paquete principal de testcontainers expone opciones genéricas
	// como WithAdditionalWaitStrategy que aplican a cualquier contenedor.
	testcontainers "github.com/testcontainers/testcontainers-go"

	// El módulo de postgres de testcontainers nos da una API de alto nivel
	// para configurar y levantar un contenedor PostgreSQL listo para usar.
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	// Wait strategies: le dicen a testcontainers cuándo el contenedor
	// está realmente listo para recibir conexiones (no solo "corriendo").
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/pav-dev98/pm-auth-svc/internal/domain"
	"github.com/pav-dev98/pm-auth-svc/internal/infrastructure/persistence/postgress"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB levanta un contenedor PostgreSQL y devuelve un repositorio
// listo para usar. Registra el cleanup automático con t.Cleanup, por lo que
// el contenedor se destruye al terminar el test sin que el caller deba hacerlo.
func setupTestDB(t *testing.T) *postgress.AuthRepository {
	t.Helper() // Marca esta función como helper: si falla, el error apunta al test que la llamó.

	ctx := context.Background()

	// Levantamos el contenedor PostgreSQL.
	// tcpostgres.Run es la API de alto nivel del módulo postgres de testcontainers.
	// "postgres:16-alpine" es la imagen Docker que se usará (Alpine = imagen liviana).
	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),

		// Wait strategy: testcontainers necesita saber cuándo el contenedor está listo.
		// WithAdditionalWaitStrategy se suma a la estrategia por defecto del módulo postgres.
		// Aquí esperamos a que el log aparezca 2 veces: PostgreSQL lo imprime una vez
		// al arrancar y otra vez al completar la inicialización. Esperar la segunda
		// ocurrencia garantiza que la DB está completamente lista para aceptar conexiones.
		// WithStartupTimeout define cuánto tiempo máximo esperamos antes de fallar.
		testcontainers.WithAdditionalWaitStrategy(
			// La palabra database system is ready to accept connections se puede mejorar la salida
			// de postgress para que sea más específica
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15 * time.Second),
		),
	)

	// require.NoError (a diferencia de assert.NoError) detiene el test inmediatamente
	// si hay error. Usamos require aquí porque sin el contenedor no tiene sentido continuar.
	require.NoError(t, err, "failed to start postgres container")

	// t.Cleanup registra una función que se ejecuta automáticamente cuando el test termina,
	// ya sea que haya pasado, fallado o llamado t.Fatal. Esto garantiza que el contenedor
	// siempre se destruya y no queden recursos huérfanos.
	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	// ConnectionString devuelve el DSN completo con el host y puerto dinámico
	// que Docker asignó al contenedor (ej: "postgres://testuser:testpass@localhost:54321/testdb").
	// El puerto es dinámico porque testcontainers mapea el 5432 interno a un puerto
	// libre del host para evitar conflictos con otras instancias.
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err, "failed to get connection string")

	// NewAuthRepository conecta a la DB y corre AutoMigrate, creando la tabla.
	repo, err := postgress.NewAuthRepository(dsn)
	require.NoError(t, err, "failed to create repository")

	return repo
}

func TestAuthRepository_Create(t *testing.T) {
	// A diferencia de los tests unitarios, los de integración NO corren en paralelo
	// entre sí cuando comparten un contenedor. Aquí cada test tiene su propio
	// contenedor (vía setupTestDB), así que sí podemos habilitarlo.
	t.Parallel()

	repo := setupTestDB(t)

	cred := domain.NewAuthCredential("create@example.com", "hashed-password")

	err := repo.Create(cred)

	// assert.NoError no detiene el test si falla (a diferencia de require).
	// Lo usamos cuando queremos que el test siga ejecutándose para ver todos los fallos.
	assert.NoError(t, err)
}

func TestAuthRepository_Create_DuplicateEmail(t *testing.T) {
	t.Parallel()

	repo := setupTestDB(t)

	cred := domain.NewAuthCredential("duplicate@example.com", "hashed-password")

	// Primera inserción: debe funcionar.
	err := repo.Create(cred)
	require.NoError(t, err)

	// Segunda inserción con el mismo email: debe fallar porque el campo
	// email tiene un uniqueIndex en GORM (ver auth_repository.go).
	duplicate := domain.NewAuthCredential("duplicate@example.com", "other-password")
	err = repo.Create(duplicate)

	assert.Error(t, err, "expected error inserting duplicate email")
}

func TestAuthRepository_FindByEmail_Found(t *testing.T) {
	t.Parallel()

	repo := setupTestDB(t)

	// Primero creamos el registro que luego vamos a buscar.
	original := domain.NewAuthCredential("find@example.com", "hashed-password")
	require.NoError(t, repo.Create(original))

	found, err := repo.FindByEmail("find@example.com")

	require.NoError(t, err)
	// require.NotNil detiene el test si found es nil, evitando un nil pointer panic
	// en los asserts siguientes.
	require.NotNil(t, found)
	assert.Equal(t, "find@example.com", found.Email)
	assert.Equal(t, "hashed-password", found.Password)
	assert.Equal(t, "user", found.Role)
	assert.True(t, found.IsActive)
}

func TestAuthRepository_FindByEmail_NotFound(t *testing.T) {
	t.Parallel()
	
	repo := setupTestDB(t)

	// Buscamos un email que nunca fue insertado.
	// El repositorio debe retornar ErrNotFound (definido en domain/errors.go),
	// no un error genérico de la DB. Esto valida que la capa de persistencia
	// traduce correctamente los errores de GORM a errores de dominio.
	result, err := repo.FindByEmail("nonexistent@example.com")

	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrNotFound,
		"repository must translate gorm.ErrRecordNotFound to domain.ErrNotFound")
}
