#!/usr/bin/env node
"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const inquirer_1 = __importDefault(require("inquirer"));
const axios_1 = __importDefault(require("axios"));
const boxen_1 = __importStar(require("boxen"));
const chalk_1 = __importDefault(require("chalk"));
const ora_1 = __importDefault(require("ora"));
const configstore_1 = __importDefault(require("configstore"));
const log = console.log;
const spinner = ora_1.default("Loading unicorns");
spinner.color = "green";
spinner.text = "Loading ...";
const url = "http://r.dsckiit.gq";
const config = new configstore_1.default("divert-cli");
log(boxen_1.default(chalk_1.default.blue(`
██████╗ ██╗██╗   ██╗███████╗██████╗ ████████╗
██╔══██╗██║██║   ██║██╔════╝██╔══██╗╚══██╔══╝
██║  ██║██║██║   ██║█████╗  ██████╔╝   ██║   
██║  ██║██║╚██╗ ██╔╝██╔══╝  ██╔══██╗   ██║   
██████╔╝██║ ╚████╔╝ ███████╗██║  ██║   ██║   
╚═════╝ ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝   ╚═╝   
`), { padding: 1, margin: 1, borderStyle: "double" }));
const main = () => {
    inquirer_1.default
        .prompt([
        {
            type: "list",
            message: "What would you like to do?",
            name: "option",
            choices: [
                "Add a new URL",
                "View list of shortened URLs",
                "Update a shortened URL",
                "Delete a shortened URL",
            ],
        },
    ])
        .then((answers) => {
        if (answers.option === "Add a new URL") {
            inquirer_1.default
                .prompt([
                {
                    type: "input",
                    name: "original_url",
                    message: "What's the URL you'd like to shorten?",
                },
                {
                    type: "input",
                    name: "shortened_url_code",
                    message: "What code would you like to assign?",
                },
            ])
                .then((answers) => {
                spinner.start();
                axios_1.default
                    .post(`${url}/api/createURL`, answers, {
                    headers: {
                        "x-auth-token": config.get("token"),
                    },
                })
                    .then((resp) => {
                    spinner.succeed();
                    log(chalk_1.default.greenBright(`Saved in DB with id ${resp.data.message}`));
                })
                    .catch((err) => log(chalk_1.default.redBright(`
          We encountered the following error:  ${err}. 
          Please try again later!`)));
            });
        }
        else if (answers.option === "View list of shortened URLs") {
            spinner.start();
            axios_1.default
                .get(`${url}/api/getAllURL`, {
                headers: {
                    "x-auth-token": config.get("token"),
                },
            })
                .then((resp) => {
                spinner.succeed();
                console.table(resp.data, [
                    "shortened_url_code",
                    "original_url",
                    "click_count",
                ]);
            })
                .catch((err) => log(chalk_1.default.redBright(`
      We encountered the following error:  ${err}. 
      Please try again later!`)));
        }
        else if (answers.option === "Update a shortened URL") {
            spinner.start();
            axios_1.default
                .get(`${url}/api/getAllURL`, {
                headers: {
                    "x-auth-token": config.get("token"),
                },
            })
                .then((resp) => {
                spinner.stop();
                console.table(resp.data, [
                    "_id",
                    "shortened_url_code",
                    "original_url",
                ]);
                inquirer_1.default
                    .prompt([
                    {
                        type: "input",
                        name: "_id",
                        message: "From the above list, enter the id of the URL you want to update",
                    },
                    {
                        type: "input",
                        name: "original_url",
                        message: "Enter new url, or the same one if there are no changes",
                    },
                    {
                        type: "input",
                        name: "shortened_url_code",
                        message: "Enter new short code, or the same one if there are no changes",
                    },
                ])
                    .then((answers) => {
                    spinner.start();
                    axios_1.default
                        .post(`${url}/api/updateURL`, answers, {
                        headers: {
                            "x-auth-token": config.get("token"),
                        },
                    })
                        .then(() => {
                        spinner.succeed();
                        log(chalk_1.default.greenBright("URL has been updated"));
                    })
                        .catch((err) => log(chalk_1.default.redBright(`
              We encountered the following error:  ${err}. 
              Please try again later!`)));
                });
            })
                .catch((err) => log(chalk_1.default.redBright(`
      We encountered the following error:  ${err}. 
      Please try again later!`)));
        }
        else if (answers.option === "Delete a shortened URL") {
            spinner.start();
            axios_1.default
                .get(`${url}/api/getAllURL`, {
                headers: {
                    "x-auth-token": config.get("token"),
                },
            })
                .then((resp) => {
                spinner.stop();
                console.table(resp.data, [
                    "_id",
                    "shortened_url_code",
                    "original_url",
                ]);
                inquirer_1.default
                    .prompt({
                    type: "input",
                    name: "_id",
                    message: "From the above list, enter the id of the URL you want to delete",
                })
                    .then((answers) => {
                    spinner.start();
                    axios_1.default
                        .post(`${url}/api/deleteURL`, answers, {
                        headers: {
                            "x-auth-token": config.get("token"),
                        },
                    })
                        .then(() => {
                        spinner.succeed();
                        log(chalk_1.default.greenBright("URL has been deleted from DB"));
                    })
                        .catch((err) => log(chalk_1.default.redBright(`
              We encountered the following error:  ${err}. 
              Please try again later!`)));
                });
            });
        }
    })
        .catch((err) => log(chalk_1.default.redBright(`
  We encountered the following error:  ${err}. 
  Please try again later!`)));
};
if (config.get("token") === undefined) {
    log("You need to be authenticated to use Divert! Enter your details below");
    inquirer_1.default
        .prompt([
        {
            type: "input",
            name: "username",
            message: "Username:",
        },
        {
            type: "password",
            name: "password",
            message: "Password:",
        },
    ])
        .then((answers) => {
        spinner.start();
        axios_1.default
            .post(`${url}/api/login`, answers)
            .then((resp) => {
            spinner.succeed();
            if (resp.data.status === "error") {
                log(chalk_1.default.redBright(`${resp.data.message}`));
                process.exit();
            }
            else {
                log(chalk_1.default.greenBright(`Authenticated!`));
                config.set({ token: resp.data.token });
            }
        })
            .then(() => main())
            .catch((err) => log(chalk_1.default.redBright(`
          We encountered the following error:  ${err}. 
          Please try again later!`)));
    });
}
else {
    main();
}
//# sourceMappingURL=index.js.map