#!/usr/bin/env node
const { spawn } = require("child_process");
const path = require("path");

const platform = process.platform;
let binary;

if (platform === "win32") binary = "koi.exe";
else if (platform === "darwin") binary = "koi-macos";
else binary = "koi-linux";

const binaryPath = path.join(__dirname, binary);

const child = spawn(binaryPath, process.argv.slice(2), { stdio: "inherit" });
child.on("exit", (code) => process.exit(code));
