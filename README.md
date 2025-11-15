# Sistema de Gestión de Biblioteca Universitaria

Sistema de gestión de biblioteca desarrollado en Go con base de datos Oracle.

## Características

- Autenticación con JWT y control de roles
- Gestión completa de libros y préstamos
- Panel administrativo con reportes
- Auditoría de todas las acciones

## Tecnologías

- **Backend**: Go + Gin Framework
- **Base de Datos**: Oracle Database
- **Autenticación**: JWT + bcrypt

## Configuración Rápida

1. **Configurar variables de entorno** (archivo `.env`):
```env
DB_USER=tu_usuario
DB_PASSWORD=tu_contraseña
DB_HOST=localhost
DB_PORT=1521
DB_SERVICE=tu_servicio
JWT_SECRET=tu_secreto_jwt
```

2. **Instalar dependencias**:
```bash
go mod download
```

3. **Ejecutar el servidor**:
```bash
go run server/main.go
```

## Credenciales por Defecto

Actualizar contraseñas con: `go run scripts/setup_users.go`

- **Admin**: admin@biblioteca.edu / admin123
- **Estudiante**: juan.perez@estudiante.edu / estudiante123
- **Profesor**: maria.lopez@profesor.edu / profesor123

## Servidor

El servidor corre en `http://localhost:8080`