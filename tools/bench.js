const spawn = require("child_process").spawn;
const fs = require("fs");
const path = require("path");
const readline = require("readline");

const CLASSES = ["druid", "hunter", "paladin", "rogue", "mage", "priest", "shaman", "warlock", "warrior"];

/**
 * @typedef SpecResult
 * @prop {string} name The class/spec string used.
 * @prop {number} timeSec
 * @prop {number} count
 * @prop {{iterations: number, nsperop: number}[]} runs The individual run results.
 * @prop {number} avg The average ns/op
 * @prop {number} dev The max. deviation from the average as a decimal.
 * @prop {number} stdDev The standard deviation.
 * 
 * @typedef BenchResult
 * @prop {string} label Label to save the bench with.
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
            }
        }
    }

    return found;
}

/**
 * Run a spec benchmark.
 * @param {string} specExt The class/spec directory string.
 * @param {number} timeSec Time in seconds for each run.
 * @param {number} count The number of runs.
 * @returns {Promise<SpecResult>} Promise resolving to the result of the bench.
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

    const specResult = {
        name: specExt,
        timeSec: timeSec,
        count: count,
        runs: [],
    }

    let output = "";

    return new Promise((resolve, reject) => {
        run.on("error", (err) => {
            console.error(err);
            reject();
        });

        run.stdout.on("data", (data) => {
            output += data;
            const split = output.split("\n");
            if (split.length > 1) {
                for (let i = 0; i < split.length - 1; i++) {
                    console.log(split[i]);
                    if (split[i].startsWith("Benchmark")) {
                        const matches = split[i].match(/Benchmark\w+ \t([\s0-9]+)\t([\s0-9]+)ns/);
                        const iterations = parseInt(matches[1].trim());
                        const nsperop = parseInt(matches[2].trim());
                        specResult.runs.push({iterations, nsperop});
                    }
                }
                output = split[split.length - 1]
            }
        });
    
        run.stderr.on("data", (data) => console.error(data));
    
        run.on("close", (code) => {
            if (code !== 0) {
                console.error("Exit code not 0, was " + code);
                reject();
            }

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

            resolve(specResult);
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
        results: {},
        totalAvg: 0,
        devMax: 0,
        stdDevMax: 0,
    }

    for (const spec of specs) {
        console.log(`Starting benchmark ${spec}, ${timeSec}s x ${count} times...`);
        const result = await runSpecBench(spec, timeSec, count);
        console.log(`Done ${spec}: ${nsToMs(result.avg)} ms +-${Math.round(result.dev * 1000) / 10}% (σ ${nsToMs(result.stdDev)} ms)`);
        benchResult.results[result.name] = result;
        benchResult.totalAvg += result.avg;
        if (benchResult.devMax < result.dev) benchResult.devMax = result.dev;
        if (benchResult.stdDevMax < result.stdDev) benchResult.stdDevMax = result.stdDev;
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
    const choices = specsWithBenchs.map((v, i) => `${i+1}: ${v}`);
    
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
            rl.question("Enter time per spec (default 3s): ", answer => {
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

    while (choice.length == 0) {
        const chosen = await new Promise((res, rej) => {
            rl.question("Which spec to run? Seperate multiple choices with a comma.\n" + choices.join("\n") + "\nChoices (Default=all): ", answer => {
                res(answer);
            });
        });

        if (!chosen) {
            choice[0] = 0;
        } else {
            const specsChosen = chosen.split(",").map(v => parseInt(v.trim()));
            for (const specNum of specsChosen) {
                if (specNum === 0) {
                    choice = [0];
                    break;
                }
                if (specNum && specNum > 0 && specNum <= choices.length) choice.push(specNum - 1);
            }
        }
    }

    rl.close();

    const specChoices = choice.map(i => specsWithBenchs[i]);

    runBenches(label, specChoices, timeSec, count);
}

benchMenu();
