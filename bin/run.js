#!/usr/bin/env node
const { spawn } = require("child_process");
const path = require("path");
const fs = require("fs");

const platform = process.platform;
let binaryFolder, binaryName;

if (platform === "win32") {
  binaryFolder = "koi_windows_amd64_v1";
  binaryName = "koi.exe";
} else if (platform === "darwin") {
  binaryFolder = "koi_darwin_amd64_v1";
  binaryName = "koi";
} else {
  binaryFolder = "koi_linux_amd64_v1";
  binaryName = "koi";
}

const binaryPath = path.join(__dirname, binaryFolder, binaryName);

// Make sure binary exists
if (!fs.existsSync(binaryPath)) {
  console.error(`Binary not found at ${binaryPath}`);
  process.exit(1);
}

const child = spawn(binaryPath, process.argv.slice(2), { stdio: "inherit" });
child.on("exit", (code) => process.exit(code));
