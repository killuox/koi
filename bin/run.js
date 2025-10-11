#!/usr/bin/env node
const { spawn } = require("child_process");
const path = require("path");

const platform = process.platform;
let binary;

if (platform === "win32") binary = "koi_windows_amd64_v1";
else if (platform === "darwin") binary = "koi_darwin_am ";
else binary = "koi_linux_amd64_v1";

const binaryPath = path.join(__dirname, binary);

const child = spawn(binaryPath, process.argv.slice(2), { stdio: "inherit" });
child.on("exit", (code) => process.exit(code));
