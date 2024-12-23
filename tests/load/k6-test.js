import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 50,              // 50 usuarios virtuales
    thresholds: {
        http_req_duration: ['p(95)<500'],  // 95% de requests bajo 500ms
        http_req_failed: ['rate<0.01'],    // menos de 1% de errores
    },
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            rate: 500,              // 500 requests por segundo
            timeUnit: '1s',
            duration: '10s',        // duración de 30 segundos en el escenario
            preAllocatedVUs: 500,    // usuarios virtuales pre-allocated
        },
    },
};

const BASE_URL = 'http://localhost:3000/api';

export default function () {
    // Test listar productos
    let productsResponse = http.get(`${BASE_URL}/ListarProductos`);
    check(productsResponse, {
        'status is 200': (r) => r.status === 200,
        'response time < 500ms': (r) => r.timings.duration < 500,
    });

    // Test obtener inventario
    let storeID = 'ba954e3f-6242-4910-bf24-e369e1dbfb68';
    let inventoryResponse = http.get(`${BASE_URL}/stores/${storeID}/inventory`);
    check(inventoryResponse, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(1);
}
