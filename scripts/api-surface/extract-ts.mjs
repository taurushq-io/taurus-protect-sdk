// Extracts the TypeScript SDK's public service surface into the shared api-surface JSON.
//
// Uses the TypeScript compiler API rather than a regex, so private/protected members and
// overloads are classified the way tsc sees them.
import fs from "fs";
import path from "path";
import { createRequire } from "module";

const [, , sdkRoot, outPath] = process.argv;
if (!sdkRoot || !outPath) {
  console.error("usage: node extract-ts.mjs <sdk-root> <out.json>");
  process.exit(2);
}

// This script lives outside the SDK package, so resolve typescript from the SDK's own
// node_modules rather than relying on it being installed next to the script.
const ts = createRequire(path.join(path.resolve(sdkRoot), "package.json"))("typescript");

const serviceDirs = [
  path.join(sdkRoot, "src/services"),
  path.join(sdkRoot, "src/services/taurus-network"),
];

const files = serviceDirs
  .filter((d) => fs.existsSync(d))
  .flatMap((d) =>
    fs
      .readdirSync(d)
      .filter((f) => f.endsWith(".ts") && f !== "index.ts")
      .map((f) => path.join(d, f))
  );

const program = ts.createProgram(files, {
  target: ts.ScriptTarget.ES2020,
  module: ts.ModuleKind.CommonJS,
  skipLibCheck: true,
  noResolve: true,
});
const checker = program.getTypeChecker();
const services = [];

for (const file of files) {
  const source = program.getSourceFile(file);
  if (!source) continue;

  ts.forEachChild(source, (node) => {
    if (!ts.isClassDeclaration(node) || !node.name) return;
    const className = node.name.text;
    if (!className.endsWith("Service")) return;
    // BaseService is the abstract base every service extends, not a service itself.
    // Counting it made this SDK report 44 where the other three report 43.
    if (className === "BaseService") return;

    const methods = [];
    for (const member of node.members) {
      if (!ts.isMethodDeclaration(member) || !member.name) continue;
      const flags = ts.getCombinedModifierFlags(member);
      if (flags & (ts.ModifierFlags.Private | ts.ModifierFlags.Protected)) continue;

      const name = member.name.getText(source);
      if (name.startsWith("_")) continue;

      const signature = checker.signatureToString(
        checker.getSignatureFromDeclaration(member),
        member,
        ts.TypeFormatFlags.NoTruncation
      );
      const doc = ts
        .getJSDocCommentsAndTags(member)
        .map((d) => (typeof d.comment === "string" ? d.comment : ""))
        .join(" ")
        .split("\n")[0]
        .trim();

      methods.push({ name, signature: `${name}${signature}`, doc });
    }

    methods.sort((a, b) => a.name.localeCompare(b.name));
    services.push({ name: className, methods });
  });
}

services.sort((a, b) => a.name.localeCompare(b.name));
fs.writeFileSync(
  outPath,
  JSON.stringify({ sdk: "typescript", source: "src/services", services }, null, 2) + "\n"
);
console.log(`typescript: ${services.length} services`);
