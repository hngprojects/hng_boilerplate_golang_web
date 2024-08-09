const core = require('@actions/core');
const github = require('@actions/github');
const { writeFile } = require('node:fs/promises');

try {
    const envFileName = core.getInput('env_file_name') || "app.env"

    const secrets = [
        "APP_NAME",
        "APP_URL",
        "SERVER_PORT",
        "DB_NAME",
        "USERNAME",
        "REDIS_PORT",
        "REDIS_HOST",
        "REDIS_DB",
        "GOOGLE_CLIENT_ID",
        "GOOGLE_CLIENT_SECRET",
        "FACEBOOK_CLIENT_ID",
        "FACEBOOK_CLIENT_SECRET",
        "SESSION_SECRET",
        "MAIL_SERVER",
        "MAIL_USERNAME",
        "MAIL_PASSWORD",
        "MAIL_PORT",
        "MIGRATE"
      ];
      

    let envString = ""

    for (let secret of secrets){
        const value = core.getInput(secret)
        envString += `${secret}=${value}\n`
    }

    const envFile = await writeFile(envFileName, envString, { encoding: 'utf-8' })

    core.setOutput("env_file", envFile);
    
    const payload = JSON.stringify(github.context.payload, undefined, 2)
    console.log(`The event payload: ${payload}`);
} catch (error) {
    core.setFailed(error.message);
}
