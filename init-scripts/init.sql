CREATE ROLE postgres WITH LOGIN SUPERUSER PASSWORD 'root';
---------------------------------------------------------------------------------------
-- Crear esquemas
CREATE SCHEMA IF NOT EXISTS catalogos;
CREATE SCHEMA IF NOT EXISTS prueba;
---------------------------------------------------------------------------------------
-- Tabla Producto
CREATE TABLE IF NOT EXISTS catalogos.Productos (
    id UUID PRIMARY KEY,
    -- UUID para identificador único
    name VARCHAR(255) NOT NULL,
    -- Nombre del producto
    description TEXT,
    -- Descripción del producto
    category VARCHAR(100),
    -- Categoría del producto
    price DECIMAL(10, 2) NOT NULL,
    -- Precio del producto (hasta 10 dígitos, 2 decimales)
    sku VARCHAR(100) UNIQUE NOT NULL,
    -- SKU único para el producto
    --campos default para control
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    -- Estado activo/inactivo para borrado lógico
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Fecha de creación
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Fecha de última modificación
);
-- Tabla Tienda
CREATE TABLE IF NOT EXISTS catalogos.Tiendas (
    id UUID PRIMARY KEY,
    -- Identificador único para la tienda
    name VARCHAR(255) NOT NULL,
    -- Nombre de la tienda
    address TEXT,
    -- Dirección de la tienda
    phone VARCHAR(15),
    -- Número de teléfono (opcional)
    --campos default para control
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    -- Estado activo/inactivo para borrado lógico
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Fecha de creación
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Fecha de última modificación
);
-- Tabla Inventario
CREATE TABLE IF NOT EXISTS prueba.Inventarios (
    id UUID PRIMARY KEY,
    -- UUID para identificador único
    productId UUID NOT NULL REFERENCES catalogos.Productos(id) ON DELETE CASCADE,
    -- Relación con Producto
    storeId UUID NOT NULL REFERENCES catalogos.Tiendas(id) ON DELETE CASCADE,
    -- Relación con Tienda
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    -- Cantidad (no negativa)
    minStock INTEGER NOT NULL CHECK (minStock >= 0),
    -- Stock mínimo (no negativo)
    --campos default para control
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    -- Estado activo/inactivo para borrado lógico
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Fecha de creación
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- Fecha de última modificación
);
-- Tabla Movimiento
CREATE TABLE IF NOT EXISTS prueba.Movimientos (
    id UUID PRIMARY KEY,
    -- UUID para identificador único
    productId UUID NOT NULL REFERENCES catalogos.Productos(id) ON DELETE CASCADE,
    -- Relación con Producto
    sourceStoreId UUID NOT NULL REFERENCES catalogos.Tiendas(id) ON DELETE CASCADE,
    -- Tienda origen
    targetStoreId UUID NOT NULL REFERENCES catalogos.Tiendas(id) ON DELETE CASCADE,
    -- Tienda destino 
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    -- Cantidad (debe ser positiva)
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Marca de tiempo
    type VARCHAR(20) NOT NULL CHECK (type IN ('IN', 'OUT', 'TRANSFER')),
    -- Tipo (IN, OUT, TRANSFER)
    --campos default para control
    activo BOOLEAN NOT NULL DEFAULT TRUE,
    -- Estado activo/inactivo para borrado lógico
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Fecha de creación
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
    )
VALUES (
        gen_random_uuid(),
        -- Genera UUID automáticamente
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
---------------------------------------------------------------------------------------
-- Insertar tiendas en catalogos.tiendas
INSERT INTO catalogos.tiendas (
        id,
        name,
        address,
        phone,
        activo,
        created_at,
        updated_at
    )
VALUES (
        gen_random_uuid(),
        'Tienda Central',
        'Av. Reforma 555, Col. Centro, CDMX',
        '555-0123',
        true,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        gen_random_uuid(),
        'Sucursal Norte',
        'Blvd. Manuel Ávila Camacho 2000, Col. San Rafael',
        '555-0124',
        true,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        gen_random_uuid(),
        'Tienda Express Sur',
        'Calzada de Tlalpan 1234, Col. Portales',
        '555-0125',
        true,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        gen_random_uuid(),
        'Sucursal Polanco',
        'Av. Presidente Masaryk 123, Col. Polanco',
        '555-0126',
        true,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ),
    (
        gen_random_uuid(),
        'Plaza Satélite',
        'Circuito Centro Comercial 2251, Cd. Satélite',
        '555-0127',
        true,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    );
---------------------------------------------------------------------------------------
-- Primero, vamos a almacenar algunos IDs existentes en variables temporales
WITH product_ids AS (
    SELECT id
    FROM catalogos.productos
    WHERE activo = true
    LIMIT 3
), store_ids AS (
    SELECT id
    FROM catalogos.tiendas
    WHERE activo = true
    LIMIT 2
), -- Ahora creamos los registros de inventario
inventory_inserts AS (
    SELECT gen_random_uuid() as id,
        p.id as product_id,
        s.id as store_id,
        floor(random() * 100 + 20)::int as quantity,
        -- Cantidad aleatoria entre 20 y 120
        10 as min_stock -- Stock mínimo fijo de 10 unidades
    FROM (
            SELECT id
            FROM product_ids
        ) p
        CROSS JOIN (
            SELECT id
            FROM store_ids
        ) s
) -- Insertamos los registros
INSERT INTO prueba.inventarios (
        id,
        productId,
        storeId,
        quantity,
        minStock,
        activo,
        created_at,
        updated_at
    )
SELECT id,
    product_id,
    store_id,
    quantity,
    min_stock,
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM inventory_inserts;
---------------------------------------------------------------------------------------
-- Insertar movimientos de ejemplo usando productos y tiendas existentes
WITH product_ids AS (
    SELECT id
    FROM catalogos.productos
    WHERE activo = true
    LIMIT 3
), store_ids AS (
    SELECT id
    FROM catalogos.tiendas
    WHERE activo = true
    LIMIT 3
)
INSERT INTO prueba.movimientos (
        id,
        productId,
        sourceStoreId,
        targetStoreId,
        quantity,
        type,
        timestamp,
        activo,
        created_at,
        updated_at
    )
SELECT -- Movimiento tipo IN
    gen_random_uuid(),
    (
        SELECT id
        FROM product_ids OFFSET 0
        LIMIT 1
    ), (
        SELECT id
        FROM store_ids OFFSET 0
        LIMIT 1
    ), (
        SELECT id
        FROM store_ids OFFSET 0
        LIMIT 1
    ), 100, 'IN', CURRENT_TIMESTAMP,
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
UNION ALL
-- Movimiento tipo OUT
SELECT gen_random_uuid(),
    (
        SELECT id
        FROM product_ids OFFSET 1
        LIMIT 1
    ), (
        SELECT id
        FROM store_ids OFFSET 1
        LIMIT 1
    ), (
        SELECT id
        FROM store_ids OFFSET 1
        LIMIT 1
    ), 50, 'OUT', CURRENT_TIMESTAMP,
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
UNION ALL
-- Movimiento tipo TRANSFER
SELECT gen_random_uuid(),
    (
        SELECT id
        FROM product_ids OFFSET 2
        LIMIT 1
    ), (
        SELECT id
        FROM store_ids OFFSET 0
        LIMIT 1
    ), -- Tienda origen
    (
        SELECT id
        FROM store_ids OFFSET 1
        LIMIT 1
    ), -- Tienda destino
    75, 'TRANSFER', CURRENT_TIMESTAMP,
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP;
---------------------------------------------------------------------------------------------------
-- Índices para optimizar consultas frecuentes
---------------------------------------------------------------------------------------
-- Índices para productos
CREATE INDEX idx_productos_sku ON catalogos.productos(sku);
CREATE INDEX idx_productos_activo ON catalogos.productos(activo);
CREATE INDEX idx_productos_category ON catalogos.productos(category);
-- Índices para tiendas
CREATE INDEX idx_tiendas_name ON catalogos.tiendas(name);
CREATE INDEX idx_tiendas_activo ON catalogos.tiendas(activo);
-- Índices compuestos para inventarios (consultas más frecuentes)
CREATE INDEX idx_inventarios_producto_tienda ON prueba.inventarios(productId, storeId);
CREATE INDEX idx_inventarios_tienda_activo ON prueba.inventarios(storeId, activo);
CREATE INDEX idx_inventarios_stock_bajo ON prueba.inventarios(quantity, minStock)
WHERE activo = true;
-- Índices para movimientos
CREATE INDEX idx_movimientos_tipo_fecha ON prueba.movimientos(type, timestamp);
CREATE INDEX idx_movimientos_producto ON prueba.movimientos(productId);
CREATE INDEX idx_movimientos_tiendas ON prueba.movimientos(sourceStoreId, targetStoreId);
-- Funciones y triggers para mantener updated_at
---------------------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION update_updated_at_column() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';
-- Trigger para productos
CREATE TRIGGER update_productos_updated_at BEFORE
UPDATE ON catalogos.productos FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Trigger para tiendas
CREATE TRIGGER update_tiendas_updated_at BEFORE
UPDATE ON catalogos.tiendas FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Trigger para inventarios
CREATE TRIGGER update_inventarios_updated_at BEFORE
UPDATE ON prueba.inventarios FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Trigger para movimientos
CREATE TRIGGER update_movimientos_updated_at BEFORE
UPDATE ON prueba.movimientos FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- Función para transferencia de inventario con validaciones
---------------------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION transfer_inventory(
        p_product_id UUID,
        p_source_store_id UUID,
        p_target_store_id UUID,
        p_quantity INTEGER
    ) RETURNS BOOLEAN AS $$
DECLARE v_source_quantity INTEGER;
v_movement_id UUID;
BEGIN -- Verificar cantidad positiva
IF p_quantity <= 0 THEN RAISE EXCEPTION 'La cantidad debe ser positiva';
END IF;
-- Verificar que las tiendas origen y destino sean diferentes
IF p_source_store_id = p_target_store_id THEN RAISE EXCEPTION 'No se puede transferir entre la misma tienda';
END IF;
-- Verificar stock disponible
SELECT quantity INTO v_source_quantity
FROM prueba.inventarios
WHERE productId = p_product_id
    AND storeId = p_source_store_id
    AND activo = true FOR
UPDATE;
IF v_source_quantity IS NULL THEN RAISE EXCEPTION 'No existe inventario en la tienda origen';
END IF;
IF v_source_quantity < p_quantity THEN RAISE EXCEPTION 'Stock insuficiente';
END IF;
-- Reducir stock en origen
UPDATE prueba.inventarios
SET quantity = quantity - p_quantity
WHERE productId = p_product_id
    AND storeId = p_source_store_id;
-- Aumentar stock en destino
INSERT INTO prueba.inventarios (
        id,
        productId,
        storeId,
        quantity,
        minStock,
        activo
    )
VALUES (
        gen_random_uuid(),
        p_product_id,
        p_target_store_id,
        p_quantity,
        10,
        true
    ) ON CONFLICT (productId, storeId) DO
UPDATE
SET quantity = prueba.inventarios.quantity + p_quantity;
-- Registrar movimiento
INSERT INTO prueba.movimientos (
        id,
        productId,
        sourceStoreId,
        targetStoreId,
        quantity,
        type,
        timestamp,
        activo
    )
VALUES (
        gen_random_uuid(),
        p_product_id,
        p_source_store_id,
        p_target_store_id,
        p_quantity,
        'TRANSFER',
        CURRENT_TIMESTAMP,
        true
    );
RETURN TRUE;
EXCEPTION
WHEN OTHERS THEN RAISE;
END;
$$ LANGUAGE plpgsql;
-- Vista para alertas de inventario
---------------------------------------------------------------------------------------
CREATE OR REPLACE VIEW prueba.vw_inventory_alerts AS
SELECT i.productId,
    i.storeId,
    p.name as product_name,
    t.name as store_name,
    i.quantity,
    i.minStock,
    CASE
        WHEN i.quantity = 0 THEN 'SIN_STOCK'
        WHEN i.quantity <= i.minStock THEN 'STOCK_BAJO'
        ELSE 'NORMAL'
    END as alert_type
FROM prueba.inventarios i
    JOIN catalogos.productos p ON i.productId = p.id
    JOIN catalogos.tiendas t ON i.storeId = t.id
WHERE i.activo = true
    AND i.quantity <= i.minStock;
-- Constraints adicionales para integridad de datos
---------------------------------------------------------------------------------------
-- Evitar transferencias a la misma tienda
ALTER TABLE prueba.movimientos
ADD CONSTRAINT check_different_stores CHECK (
        sourceStoreId != targetStoreId
        OR type != 'TRANSFER'
    );
-- Constraint para tipos de movimiento válidos
ALTER TABLE prueba.movimientos
ADD CONSTRAINT check_movement_type CHECK (type IN ('IN', 'OUT', 'TRANSFER'));