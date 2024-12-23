---------------------------------------------------------------------------------------
-- Crear esquemas
CREATE SCHEMA IF NOT EXISTS catalogos;
CREATE SCHEMA IF NOT EXISTS prueba;
---------------------------------------------------------------------------------------
-- Tabla Producto
CREATE TABLE IF NOT EXISTS catalogos.Productos (
    id UUID PRIMARY KEY, -- UUID para identificador único
    name VARCHAR(255) NOT NULL, -- Nombre del producto
    description TEXT, -- Descripción del producto
    category VARCHAR(100), -- Categoría del producto
    price DECIMAL(10, 2) NOT NULL, -- Precio del producto (hasta 10 dígitos, 2 decimales)
    sku VARCHAR(100) UNIQUE NOT NULL, -- SKU único para el producto
	--campos default para control
	activo BOOLEAN NOT NULL DEFAULT TRUE, -- Estado activo/inactivo para borrado lógico
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Fecha de creación
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Fecha de última modificación
);

-- Tabla Tienda
CREATE TABLE IF NOT EXISTS catalogos.Tiendas (
    id UUID PRIMARY KEY, -- Identificador único para la tienda
    name VARCHAR(255) NOT NULL, -- Nombre de la tienda
    address TEXT, -- Dirección de la tienda
    phone VARCHAR(15), -- Número de teléfono (opcional)

	--campos default para control
	activo BOOLEAN NOT NULL DEFAULT TRUE, -- Estado activo/inactivo para borrado lógico
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Fecha de creación
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Fecha de última modificación
);

-- Tabla Inventario
CREATE TABLE IF NOT EXISTS prueba.Inventarios (
    id UUID PRIMARY KEY, -- UUID para identificador único
    productId UUID NOT NULL REFERENCES catalogos.Productos(id) ON DELETE CASCADE, -- Relación con Producto
	storeId UUID NOT NULL REFERENCES catalogos.Tiendas(id) ON DELETE CASCADE, -- Relación con Tienda
    quantity INTEGER NOT NULL CHECK (quantity >= 0), -- Cantidad (no negativa)
    minStock INTEGER NOT NULL CHECK (minStock >= 0), -- Stock mínimo (no negativo)
	--campos default para control
	activo BOOLEAN NOT NULL DEFAULT TRUE, -- Estado activo/inactivo para borrado lógico
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Fecha de creación
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Fecha de última modificación
);

-- Tabla Movimiento
CREATE TABLE IF NOT EXISTS prueba.Movimientos (
    id UUID PRIMARY KEY, -- UUID para identificador único
    productId UUID NOT NULL REFERENCES catalogos.Productos(id) ON DELETE CASCADE, -- Relación con Producto
    sourceStoreId UUID NOT NULL REFERENCES catalogos.Tiendas(id) ON DELETE CASCADE, -- Tienda origen
    targetStoreId UUID NOT NULL REFERENCES catalogos.Tiendas(id) ON DELETE CASCADE, -- Tienda destino 
    quantity INTEGER NOT NULL CHECK (quantity > 0), -- Cantidad (debe ser positiva)
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Marca de tiempo
    type VARCHAR(20) NOT NULL CHECK (type IN ('IN', 'OUT', 'TRANSFER')), -- Tipo (IN, OUT, TRANSFER)
	--campos default para control
	activo BOOLEAN NOT NULL DEFAULT TRUE, -- Estado activo/inactivo para borrado lógico
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Fecha de creación
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Fecha de última modificación
);

---------------------------------------------------------------------------------------
-- Insertar productos en catalogos.productos
INSERT INTO catalogos.productos (
    id,
    name, 
    description, 
    category, 
    price, 
    sku, 
    activo, 
    created_at, 
    updated_at
) VALUES 
(
    gen_random_uuid(), -- Genera UUID automáticamente
    'Laptop HP Pavilion',
    'Laptop HP Pavilion con procesador Intel i5, 8GB RAM, 256GB SSD',
    'Electrónicos',
    12999.99,
    'LAP-HP-001',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Monitor Dell 27"',
    'Monitor Dell de 27 pulgadas, Full HD, 75Hz',
    'Electrónicos',
    4599.99,
    'MON-DELL-001',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Teclado Mecánico Logitech',
    'Teclado mecánico gaming RGB',
    'Periféricos',
    1299.99,
    'TEC-LOG-001',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Mouse Gaming Razer',
    'Mouse óptico gaming con 6 botones programables',
    'Periféricos',
    899.99,
    'MOU-RAZ-001',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Auriculares Sony',
    'Auriculares inalámbricos con cancelación de ruido',
    'Audio',
    2499.99,
    'AUR-SON-001',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
),
(
    gen_random_uuid(),
    'Tablet Samsung',
    'Tablet Samsung Galaxy Tab A8 10.5"',
    'Electrónicos',
    4999.99,
    'TAB-SAM-001',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);