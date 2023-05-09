import fetch from "node-fetch";
import fs from "fs";

async function formatVertical(url) {
  const response = await fetch(url);
  const t = await response.text();

  // Remove comments
  const code = t.replace(/\/\/.*|\/\*[\s\S]*?\*\//g, "");


  // Replace all spaces with newline characters, including spaces inside strings
  let insideString = false;
  let compactCode = "";
  for (let i = 0; i < code.length; i++) {
    const char = code[i];
    if (char === '"' || char === "'" || char === "`") {
      insideString = !insideString;
    }
    if (insideString) {
      if (char === "\n") {
        compactCode += " ";
      } else {
        compactCode += char;
      }
    } else {
      if (char === " ") {
        compactCode += "\n";
      } else {
        compactCode += char;
      }
    }
  }

  // Split the code into an array of tokens
  const tokenRegex = /([a-zA-Z_$][a-zA-Z_$0-9]*)|([\{\}\(\)\[\];,])|(\".*?\")|('.*?')|(`[\s\S]*?`) |(\d*\.\d+|\d+)|(\S)/g;
  const tokens = compactCode.match(tokenRegex);

  // Process each token
  let resultingCode = "";
  for (let i = 0; i < tokens.length; i++) {
    const token = tokens[i];
    if (token === "{") {
      resultingCode += "\n" + token + "\n";
    } else if (token === "}") {
      resultingCode += "\n" + token + "\n";
    } else if (token.endsWith(";")) {
      resultingCode += token + "\n";
    } else if (
      token === "function" ||
      token === "if" ||
      token === "else" ||
      token === "for" ||
      token === "while" ||
      token === "return" ||
      token === "true" ||
      token === "false"
    ) {
      resultingCode += "\n" + token + " ";
    } else if (token.startsWith('"') || token.startsWith("'") || token.startsWith("`")) {
      if (token.includes("\n")) {
        resultingCode += "\n" + token + "\n";
      } else {
        resultingCode += " " + token + " ";
      }
    } else if (
      i + 1 < tokens.length &&
      (token === "===" || token === "&&" || token === "||") &&
      tokens[i + 1][0] !== ";" &&
      tokens[i + 1][0] !== "," &&
      tokens[i + 1][0] !== ")"
    ) {
      resultingCode += " " + token + " ";
    } else {
      resultingCode += "\n" + token;
    }
  }

  return resultingCode;
}

async function formatAndSave(url, outputFile) {
  const formattedCode = await formatVertical(url);
  fs.writeFile(outputFile, formattedCode, function (err) {
    if (err) {
      console.error("Error writing file:", err);
    } else {
      console.log("Formatted code saved to file:", outputFile);
    }
  });
}

formatAndSave("https://code.jquery.com/jquery-3.6.4.min.js", "formatted.js");
