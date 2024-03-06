const spawn = require("child_process").spawn;
const fs = require("fs");
const path = require("path");
const readline = require("readline");

const CLASSES = ["druid", "hunter", "paladin", "rogue", "mage", "priest", "shaman", "warlock", "warrior"];

/**
 * @typedef SpecResult
 * @prop {string} name The class/spec string used.
 * @prop {string} funcName The name of the function used.
 * @prop {{iterations: number, nsperop: number}[]} runs The individual run results.
 * @prop {number} avg The average ns/op
 * @prop {number} dev The max. deviation from the average as a decimal.
 * @prop {number} stdDev The standard deviation.
 *
 * @typedef BenchResult
 * @prop {string} label Label to save the bench with.
 * @prop {number} timeSec The time per bench run.
 * @prop {number} count The count of runs per spec.
 * @prop {{[classSpec: string]: SpecResult}} results
 * @prop {number} totalAvg The average ns/op over all specs.
 * @prop {number} devMax The max. deviation of all specs.
 * @prop {number} stdDevMax The max. standard deviation of all specs.
 */

/**
 * Turn ns value into ms value.
 * @param {number} ns The ns number.
 * @param {number} prec The precission.
 * @returns {number} The value in milliseconds.
 */
function nsToMs(ns, prec = 3) {
    const pmult = Math.pow(10, prec);
    return Math.round(ns / 1e6 * pmult) / pmult;
}

/**
 * Get all spec benchmark targets.
 * @returns {string[]} Array of "class/spec" strings, representing the directory the test files are in.
 */
function findSpecBenchs() {
    const found = [];

    for (const classStr of CLASSES) {
        const classDir = "sim/" + classStr;
        if (!fs.existsSync(classDir)) continue;
        const classDirContent = fs.readdirSync(classDir);

        for (const specDir of classDirContent) {
            const stats = fs.statSync(path.join(classDir, specDir))
            if (stats.isDirectory() && !specDir.startsWith("_")) {
                const specDirContent = fs.readdirSync(path.join(classDir, specDir));
                for (const specFile of specDirContent) {
                    if (specFile.includes("_test.go")) {
                        const specFileData = fs.readFileSync(path.join(classDir, specDir, specFile));
                        if (specFileData.includes("func Benchmark")) {
                            found.push(`${classStr}/${specDir}`);
                        }
                    }
                }
            } else if (specDir.includes("_test.go")) {
                const specFileData = fs.readFileSync(path.join(classDir, specDir));
                if (specFileData.includes("func Benchmark")) {
                    found.push(`${classStr}`);
                }
            }
        }
    }

    return found;
}

/**
 *
 * @param {string} specExt
 * @returns {SpecResult}
 */
function newSpecResult(specExt) {
    return {
        name: specExt,
        funcName: "",
        runs: [],
        avg: 0,
        dev: 0,
        stdDev: 0,
    }
}

/**
 *
 * @param {SpecResult} specResult
 */
function finishSpecResult(specResult) {
    let min = Number.MAX_SAFE_INTEGER;
    let max = 0;
    let sum = 0;

    for (const run of specResult.runs) {
        sum += run.nsperop;
        if (min > run.nsperop) min = run.nsperop;
        if (max < run.nsperop) max = run.nsperop;
    }

    specResult.avg = Math.round(sum / specResult.runs.length);
    const maxDist = Math.max(max - specResult.avg, specResult.avg - min);
    specResult.dev = maxDist / specResult.avg;

    let variance = 0;
    for (const run of specResult.runs) {
        variance += Math.pow(run.nsperop - specResult.avg, 2);
    }
    variance /= specResult.runs.length;
    specResult.stdDev = Math.sqrt(variance);
}

/**
 * Run a spec benchmark.
 * @param {string} specExt The class/spec directory string.
 * @param {number} timeSec Time in seconds for each run.
 * @param {number} count The number of runs.
 * @returns {Promise<SpecResult[]>} Promise resolving to the result of the bench.
 */
function runSpecBench(specExt, timeSec, count) {
    const run = spawn("go", [
        "test",
        "-run=notestsplx",
        "-bench=.",
        "-cpu=1",
        `-benchtime=${timeSec}s`,
        `-count=${count}`,
        "--tags=with_db",
        `./sim/${specExt}/...`,
    ]);

    run.stdout.setEncoding("utf-8");
    run.stderr.setEncoding("utf-8");

    const specResults = [newSpecResult(specExt)];
    let specResult = specResults[0];

    let output = "";

    return new Promise((resolve, reject) => {
        run.on("error", err => {
            console.error(err);
            reject();
        });

        run.stdout.on("data", data => {
            output += data;
            const split = output.split("\n");
            if (split.length > 1) {
                for (let i = 0; i < split.length - 1; i++) {
                    console.log(split[i]);
                    if (split[i].startsWith("Benchmark")) {
                        const matches = split[i].match(/(Benchmark\w+)\s+(\d+)\s+(\d+) ns/);
                        if (specResult.funcName != "" && matches[1] != specResult.funcName) {
                            finishSpecResult(specResult);
                            specResult = newSpecResult(specExt);
                            specResults.push(specResult);
                        }
                        specResult.funcName = matches[1];
                        const iterations = parseInt(matches[2].trim());
                        const nsperop = parseInt(matches[3].trim());
                        specResult.runs.push({iterations, nsperop});
                    }
                }
                output = split[split.length - 1]
            }
        });

        run.stderr.on("data", data => console.error(data));

        run.on("close", code => {
            if (code !== 0) {
                console.error("Exit code not 0, was " + code);
                reject();
            }
            finishSpecResult(specResult)
            resolve(specResults);
        });
    });
}

/**
 * Run benchmarks and save results to file.
 * @param {string} label Label to save the bench with.
 * @param {string[]} specs The class/spec benchmarks to run.
 * @param {number} timeSec The time for each individual run.
 * @param {number} count The run count for each spec.
 */
async function runBenches(label, specs, timeSec, count) {
    /** @type {BenchResult} */
    const benchResult = {
        label: label,
        timeSec: timeSec,
        count: count,
        results: {},
        totalAvg: 0,
        devMax: 0,
        stdDevMax: 0,
    }

    for (const spec of specs) {
        console.log(`Starting benchmark ${spec}, ${timeSec}s x ${count} times...`);
        try {
            const results = await runSpecBench(spec, timeSec, count);
            for (const result of results) {
                console.log(`Done ${spec} ${result.funcName}: ${nsToMs(result.avg)} ms +-${Math.round(result.dev * 1000) / 10}% (σ ${nsToMs(result.stdDev)} ms)`);
                benchResult.results[`${result.name}_${result.funcName}`] = result;
                benchResult.totalAvg += result.avg;
                if (benchResult.devMax < result.dev) benchResult.devMax = result.dev;
                if (benchResult.stdDevMax < result.stdDev) benchResult.stdDevMax = result.stdDev;
            }
        } catch (error) {
            console.error(error);
            console.error("Skipping class/spec " + spec + " due to error!");
        }
    }

    benchResult.totalAvg = Math.round(benchResult.totalAvg / specs.length);

    const outfile = `benchres/${label}.json`;
    fs.writeFileSync(outfile, JSON.stringify(benchResult, null, 4));
    console.log("Saved result data to " + outfile);
}

/**
 * Show menu to run benchmarks.
 */
async function benchMenu() {

    const specsWithBenchs = findSpecBenchs();
    console.log("Found benchmarks: " + specsWithBenchs.join(", "));

    const rl = readline.createInterface({ input: process.stdin, output: process.stdout });

    let label = "";
    let timeSec = 0;
    let count = 0;
    let choice = [];

    while (!label) {
        label = await new Promise((res, rej) => {
            rl.question("Enter label for this run: ", answer => {
                res(answer);
            });
        });
    }

    if (label.includes("{time}")) {
        const now = new Date();
        label = label.replace("{time}", `${now.getUTCFullYear()}-${now.getUTCMonth()+1}-${now.getUTCDate()}_${now.getUTCHours()}_${now.getUTCMinutes()}`);
    }

    console.log("Using label " + label);

    while (!count) {
        const n = await new Promise((res, rej) => {
            rl.question("Enter runs per spec (default 5): ", answer => {
                res(answer);
            });
        });
        if (!n) {
            count = 5;
            break;
        }
        if (parseInt(n)) count = parseInt(n);
    }
    while (!timeSec) {
        const n = await new Promise((res, rej) => {
            rl.question("Enter time per spec (default 5s): ", answer => {
                res(answer);
            });
        });
        if (!n) {
            timeSec = 3;
            break;
        }
        if (parseInt(n)) timeSec = parseInt(n);
    }
    console.log(`Will run each spec ${timeSec}s x ${count} times.`);

    const choices = specsWithBenchs.map((v, i) => `${i+1}: ${v}`);
    while (choice.length == 0) {
        const chosen = await new Promise((res, rej) => {
            rl.question("Which spec to run? Seperate multiple choices with a comma.\n" + choices.join("\n") + "\nChoices (Default=all): ", answer => {
                res(answer);
            });
        });

        if (!chosen) {
            choice[0] = -1;
        } else {
            const specsChosen = chosen.split(",").map(v => parseInt(v.trim()));
            for (const specNum of specsChosen) {
                if (specNum === 0) {
                    choice = [-1];
                    break;
                }
                if (specNum && specNum > 0 && specNum <= choices.length) choice.push(specNum - 1);
            }
        }
    }
    if (choice[0] == -1) choice = specsWithBenchs.map((_, i) => i);

    rl.close();

    const specChoices = choice.map(i => specsWithBenchs[i]);

    runBenches(label, specChoices, timeSec, count);
}

/**
 *
 * @param {string[][]} data
 */
function comparisonStringTable(data) {
    let maxLens = [];
    for (const line of data) {
        for (let i = 0; i < line.length; i++) {
            if (!maxLens[i] || line[i].length > maxLens[i]) maxLens[i] = line[i].length;
        }
    }

    let str = "";

    for (let i = 0; i < data.length; i++) {
        str += data[i].map((v, i) => {
            if (typeof v !== "string") v = v.toString();
            return v.padEnd(maxLens[i] + 1, " ");
        }).join("\t") + "\n";
    }

    return str;
}

/**
 *
 * @param {string} resOldFile
 * @param {string} resNewFile
 */
function compareResults(resOldFile, resNewFile) {
    /** @type BenchResult */
    let oldData;
    /** @type BenchResult */
    let newData;
    try {
        let data = fs.readFileSync(path.join("benchres", resOldFile), "utf-8");
        oldData = JSON.parse(data);

        data = fs.readFileSync(path.join("benchres", resNewFile), "utf-8");
        newData = JSON.parse(data);
    } catch (error) {
        console.error(error);
    }

	const infoheader = `Comparing
Old: ${oldData.label} (Ran ${oldData.count} * ${oldData.timeSec}s per spec bench)
New: ${newData.label} (Ran ${newData.count} * ${newData.timeSec}s per spec bench)\n\n`;

    /** @type string[][] */
    const tableData = [
        ["Benchmark", "Old ms/op", "Old dev", "New ms/op", "New dev", "Delta Pct", "Delta ms"],
    ]

    for (const spec in oldData.results) {
        if (!newData.results[spec]) {
            console.log(`Spec results for ${spec} do not exist in new results.`);
            continue;
        }
        const oldSpec = oldData.results[spec];
        const newSpec = newData.results[spec];

        tableData.push([
            spec,
            nsToMs(oldSpec.avg),
            `+-${Math.round(oldSpec.dev * 1000) / 10}% (σ ${nsToMs(oldSpec.stdDev)} ms)`,
            nsToMs(newSpec.avg),
            `+-${Math.round(newSpec.dev * 1000) / 10}% (σ ${nsToMs(newSpec.stdDev)} ms)`,
            `${Math.round((newSpec.avg - oldSpec.avg) / oldSpec.avg * 1000) / 10}%`,
            nsToMs(newSpec.avg - oldSpec.avg),
        ]);
    }

    const dev = newData.totalAvg - oldData.totalAvg;
    tableData.push([
        "Total",
        nsToMs(oldData.totalAvg),
        `+-${Math.round(oldData.devMax * 1000) / 10}%`,
        nsToMs(newData.totalAvg),
        `+-${Math.round(newData.devMax * 1000) / 10}%`,
        `${Math.round(dev / oldData.totalAvg * 1000) / 100}%`,
        nsToMs(dev),
    ]);

    const csv = tableData.map(v => v.join(",")).join("\n");
    const strTable = infoheader + comparisonStringTable(tableData);
    const fname = `comp_${resOldFile.split(".")[0]}_${resNewFile.split(".")[0]}`;

	let filePath = path.join("benchres", fname+".csv");
    fs.writeFileSync(filePath, csv);
	console.log(`Wrote files ${filePath}`);

	filePath = path.join("benchres", fname+".txt");
    fs.writeFileSync(filePath, strTable);
	console.log(`Wrote files ${filePath}`);
}

async function compareMenu() {
    const rl = readline.createInterface({ input: process.stdin, output: process.stdout });

    const jsons = [];

    try {
        const files = fs.readdirSync("benchres");
        for (const file of files) {
            if (file.endsWith(".json")) {
                jsons.push(file);
            }
        }
    } catch (error) {
        console.error(error);
        return;
    }

    if (jsons.length == 0) {
        console.log("No result files found!");
        return;
    }

    const choices = jsons.map((v, i) => `${i+1}: ${v}`);
    let choice = [];

    while (choice == 0) {
        const chosen = await new Promise((res, rej) => {
            rl.question("Which results to compare? Specify old and new result seperated by a comma:\n" + choices.join("\n") + "\nChoices: ", answer => {
                res(answer);
            });
        });

        if (chosen)  {
            const split = chosen.split(",").map(v => parseInt(v.trim()));
            if (split.length != 2) continue;
            for (const num of split) {
                if (num && num >= 0 && num <= choices.length) choice.push(num - 1);
            }
            if (choice.length != 2) choice = [];
        }
    }

    rl.close();

    console.log(`Comparing ${jsons[choice[0]]} and ${jsons[choice[1]]}.`);
    compareResults(jsons[choice[0]], jsons[choice[1]]);
}

let doComp = false;
for (const arg of process.argv) {
    if (arg == "comp") {
        doComp = true;
        break;
    }
}

if (doComp) {
    compareMenu();
} else {
    benchMenu();
}
