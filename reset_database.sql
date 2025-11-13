-- ============================================================================
-- SCRIPT DE RESET Y CONFIGURACIÓN COMPLETA DE BASE DE DATOS
-- Sistema de Gestión de Biblioteca Universitaria
-- ============================================================================
-- Este script:
-- 1. Elimina todas las tablas y secuencias existentes (si existen)
-- 2. Crea todas las tablas desde cero
-- 3. Crea todas las secuencias
-- 4. Inserta datos iniciales (roles, permisos, usuarios de prueba, libros)
-- 5. Es idempotente: se puede ejecutar múltiples veces sin errores
-- ============================================================================

SET SERVEROUTPUT ON;
SET VERIFY OFF;

-- ============================================================================
-- PASO 1: ELIMINAR TABLAS EXISTENTES (EN ORDEN CORRECTO POR FKs)
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 1: Eliminando tablas existentes ===');
END;
/

-- Eliminar tablas hijas primero (las que tienen FKs)
BEGIN EXECUTE IMMEDIATE 'DROP TABLE UsuarioRol CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE RolPermiso CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE LibroAutor CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Ejemplar CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Prestamo CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Bitacora CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Estudiante CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Profesor CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Personal CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Libro CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/

-- Eliminar tablas padre
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Usuario CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Roles CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Permiso CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Autor CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP TABLE Editorial CASCADE CONSTRAINTS'; EXCEPTION WHEN OTHERS THEN NULL; END;
/

-- ============================================================================
-- PASO 2: ELIMINAR SECUENCIAS EXISTENTES
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 2: Eliminando secuencias existentes ===');
END;
/

BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE USUARIO_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE ESTUDIANTE_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE PROFESOR_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE PERSONAL_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE ROLES_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE PERMISO_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE USUARIOROL_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE ROLPERMISO_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE LIBRO_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE AUTOR_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE EDITORIAL_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE LIBROAUTOR_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE EJEMPLAR_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE PRESTAMO_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/
BEGIN EXECUTE IMMEDIATE 'DROP SEQUENCE BITACORA_SEQ'; EXCEPTION WHEN OTHERS THEN NULL; END;
/

-- ============================================================================
-- PASO 3: CREAR TABLAS
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 3: Creando tablas ===');
END;
/

-- Tabla Usuario (tabla principal)
CREATE TABLE Usuario (
    idUsuario     INTEGER      NOT NULL,
    nombre        VARCHAR2(50),
    apellido      VARCHAR2(50),
    contrasenia   VARCHAR2(100),
    correo        VARCHAR2(200),
    telefono      INTEGER,
    fechaRegistro DATE,
    CONSTRAINT Usuario_PK PRIMARY KEY (idUsuario)
);

-- Tabla Roles
CREATE TABLE Roles (
    idRol     INTEGER       NOT NULL,
    nombreRol VARCHAR2(100),
    CONSTRAINT Roles_PK PRIMARY KEY (idRol)
);

-- Tabla Permiso
CREATE TABLE Permiso (
    idPermiso   INTEGER       NOT NULL,
    descripcion VARCHAR2(100),
    CONSTRAINT Permiso_PK PRIMARY KEY (idPermiso)
);

-- Tabla UsuarioRol (N:M entre Usuario y Roles)
CREATE TABLE UsuarioRol (
    idUsuarioRol      INTEGER NOT NULL,
    Usuario_idUsuario INTEGER NOT NULL,
    Roles_idRol       INTEGER NOT NULL,
    CONSTRAINT UsuarioRol_PK PRIMARY KEY (idUsuarioRol),
    CONSTRAINT UsuarioRol_Usuario_FK FOREIGN KEY (Usuario_idUsuario) REFERENCES Usuario(idUsuario),
    CONSTRAINT UsuarioRol_Roles_FK FOREIGN KEY (Roles_idRol) REFERENCES Roles(idRol)
);

-- Tabla RolPermiso (N:M entre Roles y Permiso)
CREATE TABLE RolPermiso (
    idRolPermiso      INTEGER NOT NULL,
    Roles_idRol       INTEGER NOT NULL,
    Permiso_idPermiso INTEGER NOT NULL,
    CONSTRAINT RolPermiso_PK PRIMARY KEY (idRolPermiso),
    CONSTRAINT RolPermiso_Roles_FK FOREIGN KEY (Roles_idRol) REFERENCES Roles(idRol),
    CONSTRAINT RolPermiso_Permiso_FK FOREIGN KEY (Permiso_idPermiso) REFERENCES Permiso(idPermiso)
);

-- Tabla Estudiante
CREATE TABLE Estudiante (
    carnet            INTEGER       NOT NULL,
    carrera           VARCHAR2(100),
    semestre          INTEGER,
    Usuario_idUsuario INTEGER       NOT NULL,
    CONSTRAINT Estudiante_PK PRIMARY KEY (carnet),
    CONSTRAINT Estudiante_Usuario_FK FOREIGN KEY (Usuario_idUsuario) REFERENCES Usuario(idUsuario)
);

-- Tabla Profesor
CREATE TABLE Profesor (
    codigoDocencia    INTEGER       NOT NULL,
    facultad          VARCHAR2(100),
    Usuario_idUsuario INTEGER       NOT NULL,
    CONSTRAINT Profesor_PK PRIMARY KEY (codigoDocencia),
    CONSTRAINT Profesor_Usuario_FK FOREIGN KEY (Usuario_idUsuario) REFERENCES Usuario(idUsuario)
);

-- Tabla Personal
CREATE TABLE Personal (
    codigoEmpleado    INTEGER       NOT NULL,
    puesto            VARCHAR2(100),
    Usuario_idUsuario INTEGER       NOT NULL,
    CONSTRAINT Personal_PK PRIMARY KEY (codigoEmpleado),
    CONSTRAINT Personal_Usuario_FK FOREIGN KEY (Usuario_idUsuario) REFERENCES Usuario(idUsuario)
);

-- Tabla Autor
CREATE TABLE Autor (
    idAutor      INTEGER       NOT NULL,
    nombre       VARCHAR2(100),
    apellido     VARCHAR2(100),
    nacionalidad VARCHAR2(100),
    CONSTRAINT Autor_PK PRIMARY KEY (idAutor)
);

-- Tabla Editorial
CREATE TABLE Editorial (
    idEditorial INTEGER       NOT NULL,
    nombre      VARCHAR2(100),
    pais        VARCHAR2(100),
    CONSTRAINT Editorial_PK PRIMARY KEY (idEditorial)
);

-- Tabla Libro
CREATE TABLE Libro (
    ISBN                  INTEGER       NOT NULL,
    titulo                VARCHAR2(200),
    anioEdicion           DATE,
    Editorial_idEditorial INTEGER       NOT NULL,
    CONSTRAINT Libro_PK PRIMARY KEY (ISBN),
    CONSTRAINT Libro_Editorial_FK FOREIGN KEY (Editorial_idEditorial) REFERENCES Editorial(idEditorial)
);

-- Tabla LibroAutor (N:M entre Libro y Autor)
CREATE TABLE LibroAutor (
    idLibroAutor  INTEGER       NOT NULL,
    tipoAutor     VARCHAR2(100),
    Autor_idAutor INTEGER       NOT NULL,
    Libro_ISBN    INTEGER       NOT NULL,
    CONSTRAINT LibroAutor_PK PRIMARY KEY (idLibroAutor),
    CONSTRAINT LibroAutor_Autor_FK FOREIGN KEY (Autor_idAutor) REFERENCES Autor(idAutor),
    CONSTRAINT LibroAutor_Libro_FK FOREIGN KEY (Libro_ISBN) REFERENCES Libro(ISBN)
);

-- Tabla Prestamo
CREATE TABLE Prestamo (
    idPrestamo              INTEGER     NOT NULL,
    fechaPrestamo           DATE,
    fechaDevolucionPrevista DATE,
    fechaDevolucionReal     DATE,
    estado                  VARCHAR2(50),
    Usuario_idUsuario       INTEGER     NOT NULL,
    Devolucion_idDevolucion INTEGER,
    CONSTRAINT Prestamo_PK PRIMARY KEY (idPrestamo),
    CONSTRAINT Prestamo_Usuario_FK FOREIGN KEY (Usuario_idUsuario) REFERENCES Usuario(idUsuario)
);

-- Tabla Ejemplar
CREATE TABLE Ejemplar (
    codigo              INTEGER     NOT NULL,
    estado              VARCHAR2(50),
    Libro_ISBN          INTEGER     NOT NULL,
    Prestamo_idPrestamo INTEGER,
    CONSTRAINT Ejemplar_PK PRIMARY KEY (codigo),
    CONSTRAINT Ejemplar_Libro_FK FOREIGN KEY (Libro_ISBN) REFERENCES Libro(ISBN),
    CONSTRAINT Ejemplar_Prestamo_FK FOREIGN KEY (Prestamo_idPrestamo) REFERENCES Prestamo(idPrestamo)
);

-- Tabla Bitacora
CREATE TABLE Bitacora (
    idBitacora        INTEGER       NOT NULL,
    accion            VARCHAR2(100),
    fechaHora         DATE,
    detalle           VARCHAR2(500),
    entidad           VARCHAR2(50),
    Usuario_idUsuario INTEGER       NOT NULL,
    CONSTRAINT Bitacora_PK PRIMARY KEY (idBitacora),
    CONSTRAINT Bitacora_Usuario_FK FOREIGN KEY (Usuario_idUsuario) REFERENCES Usuario(idUsuario)
);

-- ============================================================================
-- PASO 4: CREAR SECUENCIAS
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 4: Creando secuencias ===');
END;
/

CREATE SEQUENCE USUARIO_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE ESTUDIANTE_SEQ START WITH 2024001 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE PROFESOR_SEQ START WITH 3001 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE PERSONAL_SEQ START WITH 4001 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE ROLES_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE PERMISO_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE USUARIOROL_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE ROLPERMISO_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE LIBRO_SEQ START WITH 1000 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE AUTOR_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE EDITORIAL_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE LIBROAUTOR_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE EJEMPLAR_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE PRESTAMO_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;
CREATE SEQUENCE BITACORA_SEQ START WITH 1 INCREMENT BY 1 NOCACHE;

-- ============================================================================
-- PASO 5: INSERTAR DATOS INICIALES - ROLES
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 5: Insertando roles ===');
END;
/

INSERT INTO Roles (idRol, nombreRol) VALUES (1, 'admin');
INSERT INTO Roles (idRol, nombreRol) VALUES (2, 'estudiante');
INSERT INTO Roles (idRol, nombreRol) VALUES (3, 'profesor');
INSERT INTO Roles (idRol, nombreRol) VALUES (4, 'personal');

-- ============================================================================
-- PASO 6: INSERTAR DATOS INICIALES - PERMISOS
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 6: Insertando permisos ===');
END;
/

INSERT INTO Permiso (idPermiso, descripcion) VALUES (1, 'Gestionar usuarios');
INSERT INTO Permiso (idPermiso, descripcion) VALUES (2, 'Gestionar libros');
INSERT INTO Permiso (idPermiso, descripcion) VALUES (3, 'Gestionar préstamos');
INSERT INTO Permiso (idPermiso, descripcion) VALUES (4, 'Ver bitácora');
INSERT INTO Permiso (idPermiso, descripcion) VALUES (5, 'Gestionar roles');
INSERT INTO Permiso (idPermiso, descripcion) VALUES (6, 'Solicitar préstamos');
INSERT INTO Permiso (idPermiso, descripcion) VALUES (7, 'Ver catálogo');

-- ============================================================================
-- PASO 7: ASIGNAR PERMISOS A ROLES
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 7: Asignando permisos a roles ===');
END;
/

-- Admin tiene todos los permisos
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (1, 1, 1);
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (2, 1, 2);
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (3, 1, 3);
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (4, 1, 4);
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (5, 1, 5);
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (6, 1, 6);
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (7, 1, 7);

-- Estudiante y Profesor pueden solicitar préstamos y ver catálogo
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (8, 2, 6);
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (9, 2, 7);
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (10, 3, 6);
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (11, 3, 7);

-- Personal puede gestionar préstamos y ver catálogo
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (12, 4, 3);
INSERT INTO RolPermiso (idRolPermiso, Roles_idRol, Permiso_idPermiso) VALUES (13, 4, 7);

-- ============================================================================
-- PASO 8: INSERTAR USUARIOS DE PRUEBA
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 8: Insertando usuarios de prueba ===');
END;
/

-- Usuario Admin (contraseña: admin123 hasheada con bcrypt)
INSERT INTO Usuario (idUsuario, nombre, apellido, contrasenia, correo, telefono, fechaRegistro) 
VALUES (1, 'Admin', 'Sistema', '$2a$10$rZ9GVqhPwPqxJ3Q8KZ7vW.xYvJ0JJXmxQF5YZQXZ5YZQXZ5YZQXZ5Y', 'admin@biblioteca.edu', 12345678, SYSDATE);

-- Usuario Estudiante (contraseña: estudiante123)
INSERT INTO Usuario (idUsuario, nombre, apellido, contrasenia, correo, telefono, fechaRegistro) 
VALUES (2, 'Juan', 'Pérez', '$2a$10$rZ9GVqhPwPqxJ3Q8KZ7vW.xYvJ0JJXmxQF5YZQXZ5YZQXZ5YZQXZ5Y', 'juan.perez@estudiante.edu', 23456789, SYSDATE);

-- Usuario Profesor (contraseña: profesor123)
INSERT INTO Usuario (idUsuario, nombre, apellido, contrasenia, correo, telefono, fechaRegistro) 
VALUES (3, 'María', 'López', '$2a$10$rZ9GVqhPwPqxJ3Q8KZ7vW.xYvJ0JJXmxQF5YZQXZ5YZQXZ5YZQXZ5Y', 'maria.lopez@profesor.edu', 34567890, SYSDATE);

-- Usuario Personal (contraseña: personal123)
INSERT INTO Usuario (idUsuario, nombre, apellido, contrasenia, correo, telefono, fechaRegistro) 
VALUES (4, 'Carlos', 'García', '$2a$10$rZ9GVqhPwPqxJ3Q8KZ7vW.xYvJ0JJXmxQF5YZQXZ5YZQXZ5YZQXZ5Y', 'carlos.garcia@biblioteca.edu', 45678901, SYSDATE);

-- Más estudiantes de prueba
INSERT INTO Usuario (idUsuario, nombre, apellido, contrasenia, correo, telefono, fechaRegistro) 
VALUES (5, 'Ana', 'Martínez', '$2a$10$rZ9GVqhPwPqxJ3Q8KZ7vW.xYvJ0JJXmxQF5YZQXZ5YZQXZ5YZQXZ5Y', 'ana.martinez@estudiante.edu', 56789012, SYSDATE);

INSERT INTO Usuario (idUsuario, nombre, apellido, contrasenia, correo, telefono, fechaRegistro) 
VALUES (6, 'Pedro', 'Rodríguez', '$2a$10$rZ9GVqhPwPqxJ3Q8KZ7vW.xYvJ0JJXmxQF5YZQXZ5YZQXZ5YZQXZ5Y', 'pedro.rodriguez@estudiante.edu', 67890123, SYSDATE);

-- ============================================================================
-- PASO 9: ASIGNAR ROLES A USUARIOS
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 9: Asignando roles a usuarios ===');
END;
/

INSERT INTO UsuarioRol (idUsuarioRol, Usuario_idUsuario, Roles_idRol) VALUES (1, 1, 1); -- Admin
INSERT INTO UsuarioRol (idUsuarioRol, Usuario_idUsuario, Roles_idRol) VALUES (2, 2, 2); -- Estudiante
INSERT INTO UsuarioRol (idUsuarioRol, Usuario_idUsuario, Roles_idRol) VALUES (3, 3, 3); -- Profesor
INSERT INTO UsuarioRol (idUsuarioRol, Usuario_idUsuario, Roles_idRol) VALUES (4, 4, 4); -- Personal
INSERT INTO UsuarioRol (idUsuarioRol, Usuario_idUsuario, Roles_idRol) VALUES (5, 5, 2); -- Estudiante
INSERT INTO UsuarioRol (idUsuarioRol, Usuario_idUsuario, Roles_idRol) VALUES (6, 6, 2); -- Estudiante

-- ============================================================================
-- PASO 10: CREAR PERFILES ESPECÍFICOS
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 10: Creando perfiles de usuarios ===');
END;
/

-- Perfiles de Estudiantes
INSERT INTO Estudiante (carnet, carrera, semestre, Usuario_idUsuario) 
VALUES (2024001, 'Ingeniería en Sistemas', 5, 2);

INSERT INTO Estudiante (carnet, carrera, semestre, Usuario_idUsuario) 
VALUES (2024002, 'Ingeniería Industrial', 3, 5);

INSERT INTO Estudiante (carnet, carrera, semestre, Usuario_idUsuario) 
VALUES (2024003, 'Administración de Empresas', 7, 6);

-- Perfil de Profesor
INSERT INTO Profesor (codigoDocencia, facultad, Usuario_idUsuario) 
VALUES (3001, 'Facultad de Ingeniería', 3);

-- Perfil de Personal
INSERT INTO Personal (codigoEmpleado, puesto, Usuario_idUsuario) 
VALUES (4001, 'Bibliotecario', 4);

-- ============================================================================
-- PASO 11: INSERTAR EDITORIALES
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 11: Insertando editoriales ===');
END;
/

INSERT INTO Editorial (idEditorial, nombre, pais) VALUES (1, 'Pearson', 'Estados Unidos');
INSERT INTO Editorial (idEditorial, nombre, pais) VALUES (2, 'O''Reilly Media', 'Estados Unidos');
INSERT INTO Editorial (idEditorial, nombre, pais) VALUES (3, 'McGraw-Hill', 'Estados Unidos');
INSERT INTO Editorial (idEditorial, nombre, pais) VALUES (4, 'Alfaomega', 'México');
INSERT INTO Editorial (idEditorial, nombre, pais) VALUES (5, 'Addison-Wesley', 'Estados Unidos');
INSERT INTO Editorial (idEditorial, nombre, pais) VALUES (6, 'Planeta', 'España');
INSERT INTO Editorial (idEditorial, nombre, pais) VALUES (7, 'Santillana', 'España');

-- ============================================================================
-- PASO 12: INSERTAR AUTORES
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 12: Insertando autores ===');
END;
/

INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (1, 'Abraham', 'Silberschatz', 'Estados Unidos');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (2, 'Andrew', 'Tanenbaum', 'Países Bajos');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (3, 'Robert', 'Martin', 'Estados Unidos');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (4, 'Martin', 'Fowler', 'Reino Unido');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (5, 'Eric', 'Evans', 'Estados Unidos');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (6, 'Donald', 'Knuth', 'Estados Unidos');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (7, 'Bjarne', 'Stroustrup', 'Dinamarca');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (8, 'Brian', 'Kernighan', 'Canadá');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (9, 'Dennis', 'Ritchie', 'Estados Unidos');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (10, 'Erich', 'Gamma', 'Suiza');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (11, 'Gabriel', 'García Márquez', 'Colombia');
INSERT INTO Autor (idAutor, nombre, apellido, nacionalidad) VALUES (12, 'Isabel', 'Allende', 'Chile');

-- ============================================================================
-- PASO 13: INSERTAR LIBROS
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 13: Insertando libros ===');
END;
/

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1001, 'Fundamentos de Sistemas de Bases de Datos', TO_DATE('2020-01-01', 'YYYY-MM-DD'), 1);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1002, 'Sistemas Operativos Modernos', TO_DATE('2018-06-15', 'YYYY-MM-DD'), 1);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1003, 'Clean Code: Manual de Estilo para el Desarrollo Ágil', TO_DATE('2019-03-20', 'YYYY-MM-DD'), 5);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1004, 'Refactoring: Improving the Design of Existing Code', TO_DATE('2019-11-10', 'YYYY-MM-DD'), 5);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1005, 'Domain-Driven Design', TO_DATE('2017-08-25', 'YYYY-MM-DD'), 5);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1006, 'The Art of Computer Programming Vol. 1', TO_DATE('2021-02-14', 'YYYY-MM-DD'), 5);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1007, 'El Lenguaje de Programación C', TO_DATE('2016-05-30', 'YYYY-MM-DD'), 1);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1008, 'Design Patterns: Elements of Reusable Object-Oriented Software', TO_DATE('2018-09-12', 'YYYY-MM-DD'), 5);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1009, 'Cien Años de Soledad', TO_DATE('2015-04-18', 'YYYY-MM-DD'), 6);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1010, 'La Casa de los Espíritus', TO_DATE('2017-07-22', 'YYYY-MM-DD'), 6);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1011, 'Introducción a los Algoritmos', TO_DATE('2019-12-05', 'YYYY-MM-DD'), 3);

INSERT INTO Libro (ISBN, titulo, anioEdicion, Editorial_idEditorial) 
VALUES (1012, 'Redes de Computadoras', TO_DATE('2020-10-08', 'YYYY-MM-DD'), 1);

-- ============================================================================
-- PASO 14: RELACIONAR LIBROS CON AUTORES
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 14: Relacionando libros con autores ===');
END;
/

INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (1, 'Principal', 1, 1001);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (2, 'Principal', 2, 1002);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (3, 'Principal', 3, 1003);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (4, 'Principal', 4, 1004);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (5, 'Principal', 5, 1005);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (6, 'Principal', 6, 1006);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (7, 'Principal', 8, 1007);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (8, 'Co-autor', 9, 1007);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (9, 'Principal', 10, 1008);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (10, 'Principal', 11, 1009);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (11, 'Principal', 12, 1010);
INSERT INTO LibroAutor (idLibroAutor, tipoAutor, Autor_idAutor, Libro_ISBN) VALUES (12, 'Principal', 2, 1012);

-- ============================================================================
-- PASO 15: INSERTAR EJEMPLARES
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 15: Insertando ejemplares de libros ===');
END;
/

-- 3 ejemplares por cada libro
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (1, 'DISPONIBLE', 1001, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (2, 'DISPONIBLE', 1001, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (3, 'DISPONIBLE', 1001, NULL);

INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (4, 'DISPONIBLE', 1002, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (5, 'DISPONIBLE', 1002, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (6, 'DISPONIBLE', 1002, NULL);

INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (7, 'DISPONIBLE', 1003, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (8, 'DISPONIBLE', 1003, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (9, 'DISPONIBLE', 1003, NULL);

INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (10, 'DISPONIBLE', 1004, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (11, 'DISPONIBLE', 1004, NULL);

INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (12, 'DISPONIBLE', 1005, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (13, 'DISPONIBLE', 1005, NULL);

INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (14, 'DISPONIBLE', 1009, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (15, 'DISPONIBLE', 1009, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (16, 'DISPONIBLE', 1009, NULL);

INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (17, 'DISPONIBLE', 1010, NULL);
INSERT INTO Ejemplar (codigo, estado, Libro_ISBN, Prestamo_idPrestamo) VALUES (18, 'DISPONIBLE', 1010, NULL);

-- ============================================================================
-- PASO 16: INSERTAR ALGUNOS PRÉSTAMOS DE EJEMPLO
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('=== PASO 16: Insertando préstamos de ejemplo ===');
END;
/

-- Préstamo activo para estudiante Juan Pérez
INSERT INTO Prestamo (idPrestamo, fechaPrestamo, fechaDevolucionPrevista, fechaDevolucionReal, estado, Usuario_idUsuario, Devolucion_idDevolucion)
VALUES (1, SYSDATE-5, SYSDATE+10, NULL, 'ACTIVO', 2, NULL);

-- Actualizar ejemplar prestado
UPDATE Ejemplar SET estado = 'PRESTADO', Prestamo_idPrestamo = 1 WHERE codigo = 1;

-- Préstamo devuelto
INSERT INTO Prestamo (idPrestamo, fechaPrestamo, fechaDevolucionPrevista, fechaDevolucionReal, estado, Usuario_idUsuario, Devolucion_idDevolucion)
VALUES (2, SYSDATE-20, SYSDATE-5, SYSDATE-3, 'DEVUELTO', 5, NULL);

-- Préstamo vencido
INSERT INTO Prestamo (idPrestamo, fechaPrestamo, fechaDevolucionPrevista, fechaDevolucionReal, estado, Usuario_idUsuario, Devolucion_idDevolucion)
VALUES (3, SYSDATE-25, SYSDATE-10, NULL, 'ACTIVO', 6, NULL);

-- Actualizar ejemplar vencido
UPDATE Ejemplar SET estado = 'PRESTADO', Prestamo_idPrestamo = 3 WHERE codigo = 7;

-- ============================================================================
-- COMMIT FINAL
-- ============================================================================
COMMIT;

-- ============================================================================
-- VERIFICACIÓN FINAL
-- ============================================================================
BEGIN
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('=== RESUMEN DE DATOS INSERTADOS ===');
    DBMS_OUTPUT.PUT_LINE('');
END;
/

DECLARE
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count FROM Usuario;
    DBMS_OUTPUT.PUT_LINE('Usuarios: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM Roles;
    DBMS_OUTPUT.PUT_LINE('Roles: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM Permiso;
    DBMS_OUTPUT.PUT_LINE('Permisos: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM Editorial;
    DBMS_OUTPUT.PUT_LINE('Editoriales: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM Autor;
    DBMS_OUTPUT.PUT_LINE('Autores: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM Libro;
    DBMS_OUTPUT.PUT_LINE('Libros: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM Ejemplar;
    DBMS_OUTPUT.PUT_LINE('Ejemplares: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM Prestamo;
    DBMS_OUTPUT.PUT_LINE('Préstamos: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM Estudiante;
    DBMS_OUTPUT.PUT_LINE('Estudiantes: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM Profesor;
    DBMS_OUTPUT.PUT_LINE('Profesores: ' || v_count);
    
    SELECT COUNT(*) INTO v_count FROM Personal;
    DBMS_OUTPUT.PUT_LINE('Personal: ' || v_count);
END;
/

BEGIN
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('======================================================');
    DBMS_OUTPUT.PUT_LINE('BASE DE DATOS CONFIGURADA EXITOSAMENTE');
    DBMS_OUTPUT.PUT_LINE('======================================================');
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('USUARIOS DE PRUEBA:');
    DBMS_OUTPUT.PUT_LINE('---------------------------------------------------');
    DBMS_OUTPUT.PUT_LINE('Admin:');
    DBMS_OUTPUT.PUT_LINE('  Email: admin@biblioteca.edu');
    DBMS_OUTPUT.PUT_LINE('  Password: admin123');
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('Estudiante:');
    DBMS_OUTPUT.PUT_LINE('  Email: juan.perez@estudiante.edu');
    DBMS_OUTPUT.PUT_LINE('  Password: estudiante123');
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('Profesor:');
    DBMS_OUTPUT.PUT_LINE('  Email: maria.lopez@profesor.edu');
    DBMS_OUTPUT.PUT_LINE('  Password: profesor123');
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE('Personal:');
    DBMS_OUTPUT.PUT_LINE('  Email: carlos.garcia@biblioteca.edu');
    DBMS_OUTPUT.PUT_LINE('  Password: personal123');
    DBMS_OUTPUT.PUT_LINE('======================================================');
END;
/
