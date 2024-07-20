module.exports = {
  apps: [
    {
      name: "run_production_app",
      script: "/home/vicradon/deployments/production/production_app",
      env: {
        SERVER_PORT: 9000,
        DB_NAME: "production_db",
        USERNAME: "production_user",
        APP_NAME: "production",
        APP_URL: "http://localhost:9000",
      },
    },
    {
      name: "run_staging_app",
      script: "/home/vicradon/deployments/staging/staging_app",
      env: {
        SERVER_PORT: 8000,
        DB_NAME: "staging_db",
        USERNAME: "staging_user",
        APP_NAME: "staging",
        APP_URL: "http://localhost:8000",
      },
    },
    {
      name: "run_development_app",
      script: "/home/vicradon/deployments/development/development_app",
      env: {
        SERVER_PORT: 7000,
        DB_NAME: "development_db",
        USERNAME: "development_user",
        APP_NAME: "development",
        APP_URL: "http://localhost:7000",
      },
    },
  ],
};
