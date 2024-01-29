# Use the latest keycloak image from quay.io
FROM quay.io/keycloak/keycloak:latest

# Enable health and metrics support
ENV KC_HEALTH_ENABLED=true
ENV KC_METRICS_ENABLED=true

# Configure a database vendor
ENV KC_DB=mysql

# Specify the database URL, username and password
ENV KC_DB_URL=jdbc:mysql://localhost:3306/keycloak
ENV KC_DB_USERNAME=root
ENV KC_DB_PASSWORD=sg46sg46

# Specify the hostname
ENV KC_HOSTNAME=localhost

# Start the server
ENTRYPOINT ["/opt/keycloak/bin/kc.sh"]

