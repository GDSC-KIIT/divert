#!/usr/bin/env node
import inquirer from "inquirer";
import axios from "axios";
import boxen, { BorderStyle } from "boxen";
import chalk from "chalk";
import ora from "ora";
import Configstore from "configstore";

const log = console.log;
const spinner = ora("Loading unicorns");
spinner.color = "green";
spinner.text = "Loading ...";

const url = "http://r.dsckiit.gq";
const config = new Configstore("divert-cli");

// config.delete('token');

log(
  boxen(
    chalk.blue(`
██████╗ ██╗██╗   ██╗███████╗██████╗ ████████╗
██╔══██╗██║██║   ██║██╔════╝██╔══██╗╚══██╔══╝
██║  ██║██║██║   ██║█████╗  ██████╔╝   ██║   
██║  ██║██║╚██╗ ██╔╝██╔══╝  ██╔══██╗   ██║   
██████╔╝██║ ╚████╔╝ ███████╗██║  ██║   ██║   
╚═════╝ ╚═╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝   ╚═╝   
`),
    { padding: 1, margin: 1, borderStyle: BorderStyle.Double }
  )
);

// main function contains core functionality
const main = () => {
  inquirer
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
      // add URL
      if (answers.option === "Add a new URL") {
        // get values
        inquirer
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
            // request
            spinner.start();
            axios
              .post(`${url}/api/createURL`, answers, {
                headers: {
                  "x-auth-token": config.get("token"),
                },
              })
              .then((resp) => {
                spinner.succeed();
                log(
                  chalk.greenBright(`Saved in DB with id ${resp.data.message}`)
                );
              })
              .catch((err) =>
                log(
                  chalk.redBright(`
          We encountered the following error:  ${err}. 
          Please try again later!`)
                )
              );
          });
      }
      // view URL
      else if (answers.option === "View list of shortened URLs") {
        spinner.start();
        axios
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
          .catch((err) =>
            log(
              chalk.redBright(`
      We encountered the following error:  ${err}. 
      Please try again later!`)
            )
          );
      }
      // update URL
      else if (answers.option === "Update a shortened URL") {
        // get list to show user
        spinner.start();
        axios
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
            // after showing
            inquirer
              .prompt([
                {
                  type: "input",
                  name: "_id",
                  message:
                    "From the above list, enter the id of the URL you want to update",
                },
                {
                  type: "input",
                  name: "original_url",
                  message:
                    "Enter new url, or the same one if there are no changes",
                },
                {
                  type: "input",
                  name: "shortened_url_code",
                  message:
                    "Enter new short code, or the same one if there are no changes",
                },
              ])
              // update with inputs
              .then((answers) => {
                spinner.start();
                axios
                  .post(`${url}/api/updateURL`, answers, {
                    headers: {
                      "x-auth-token": config.get("token"),
                    },
                  })
                  .then(() => {
                    spinner.succeed();
                    log(chalk.greenBright("URL has been updated"));
                  })
                  .catch((err) =>
                    log(
                      chalk.redBright(`
              We encountered the following error:  ${err}. 
              Please try again later!`)
                    )
                  );
              });
          })
          .catch((err) =>
            log(
              chalk.redBright(`
      We encountered the following error:  ${err}. 
      Please try again later!`)
            )
          );
      }
      // delete a URL
      else if (answers.option === "Delete a shortened URL") {
        // get list to show user
        spinner.start();
        axios
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

            inquirer
              .prompt({
                type: "input",
                name: "_id",
                message:
                  "From the above list, enter the id of the URL you want to delete",
              })
              .then((answers) => {
                // delete request
                spinner.start();
                axios
                  .post(`${url}/api/deleteURL`, answers, {
                    headers: {
                      "x-auth-token": config.get("token"),
                    },
                  })
                  .then(() => {
                    spinner.succeed();
                    log(chalk.greenBright("URL has been deleted from DB"));
                  })
                  .catch((err) =>
                    log(
                      chalk.redBright(`
              We encountered the following error:  ${err}. 
              Please try again later!`)
                    )
                  );
              });
          });
      }
    })
    // handle first inquire error
    .catch((err) =>
      log(
        chalk.redBright(`
  We encountered the following error:  ${err}. 
  Please try again later!`)
      )
    );
};

// auth

// if JWT isn't stored hit login
if (config.get("token") === undefined) {
  log("You need to be authenticated to use Divert! Enter your details below");
  inquirer
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
      axios
        .post(`${url}/api/login`, answers)
        .then((resp) => {
          spinner.succeed();
          log(chalk.greenBright(`Authenticated!`));
          config.set({ token: resp.data.token });
        })
        .then(() => main())
        .catch((err) =>
          log(
            chalk.redBright(`
          We encountered the following error:  ${err}. 
          Please try again later!`)
          )
        );
    });
} else {
  main();
}
