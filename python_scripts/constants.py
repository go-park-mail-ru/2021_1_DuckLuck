ALL_LOCAL_BUILD_TARGETS = [
    'api-server',
    'api-db',
    'session-service',
    'cart-service',
    'auth-service',
]

LOCAL_UP_TARGETS_WITH_BUILD = [
    'duckluckmarket/api-db',
    'duckluckmarket/api-server',
    'duckluckmarket/session-service',
    'duckluckmarket/cart-service',
    'duckluckmarket/auth-service',
]


LOCAL_UP_TARGETS_FROM_DOCKERHUB = [
    'auth-db',
    'session-db',
    'cart-db',
    'node-exporter',
    'prometheus',
    'alertmanager-bot',
    'alertmanager',
    'grafana',
]


IMAGES_FROM_DOCKERHUB = [
    'postgres',
    'redis',
    'prom/node-exporter',
    'prom/prometheus',
    'prom/alertmanager',
    'grafana/grafana:latest',
    'metalmatze/alertmanager-bot',
]