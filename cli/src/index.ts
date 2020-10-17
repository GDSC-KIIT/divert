#!/usr/bin/env node
import inquirer from "inquirer";
import axios from "axios";
import boxen, { BorderStyle } from "boxen";
import chalk from "chalk";

const log = console.log;

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
    if (answers.option === "View list of shortened URLs") {
      axios
        .get("http://dsckiit-divert.herokuapp.com/api/getAllURL")
        .then((resp) => {
          console.table(resp.data, [
            "shortened_url_code",
            "original_url",
            "click_count",
          ]);
        });
    }
  });
