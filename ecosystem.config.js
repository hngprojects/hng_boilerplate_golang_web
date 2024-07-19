module.exports = {
  apps: [
    {
      name: 'run_learnai_prod',
      script: '/home/vicradon/deployments/production/learnai_prod',
      env: {
        SERVER_PORT: 9000,
        DB_NAME: 'db_name',
        APP_NAME: 'production',
        APP_URL: 'http://localhost:9000'
      }
    },
    {
      name: 'run_learnai_staging',
      script: '/home/vicradon/deployments/staging/learnai_staging',
      env: {
        SERVER_PORT: 8000,
        DB_NAME: 'db_name',
        APP_NAME: 'staging',
        APP_URL: 'http://localhost:8000'
      }
    },
    {
      name: 'run_learnai_dev',
      script: '/home/vicradon/deployments/development/learnai_dev',
      env: {
        SERVER_PORT: 7000,
        DB_NAME: 'db_name',
        APP_NAME: 'development',
        APP_URL: 'http://localhost:7000'
      }
    }
  ]
};
