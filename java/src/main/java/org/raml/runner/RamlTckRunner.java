package org.raml.runner;

import org.raml.parsers.IParser;
import org.raml.parsers.WebApiParser;
import org.raml.parsers.RamlJavaParser;
import org.raml.utils.Utils;

import picocli.CommandLine;
import picocli.CommandLine.ParameterException;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;
import picocli.CommandLine.Parameters;
import org.json.simple.JSONObject;
import org.json.simple.JSONArray;

import java.util.*;


@Command(name = "raml-tck-runner", mixinStandardHelpOptions = true, version = "1.0.0")
public class RamlTckRunner implements Runnable {
  @Option(names = "--parser", description = "name of a parser to run")
  String parserName;

  @Option(names = "--outdir", description = "output JSON report directory")
  String outdir = "./";

  @Option(names = "--branch", description = "raml-tck branch to load RAML files from")
  String branch;

  public IParser pickParser() {
    IParser parser;
    switch (parserName) {
      case "webapi-parser":
        parser = new WebApiParser();
        break;
      case "raml-java-parser":
        parser = new RamlJavaParser();
        break;
      default:
        throw new ParameterException(
          new CommandLine(this),
          "Not supported parser: " + parserName);
    }
    return parser;
  }

  public void run() {
    IParser parser = this.pickParser();
    String exDir = Utils.cloneTckRepo(branch);
    List<String> fileList = Utils.listRamls(exDir);

    JSONObject report = new JSONObject();
    report.put("parser", parserName + "(java)");
    report.put("branch", branch);
    JSONArray results = new JSONArray();

    Boolean success;
    String error;
    JSONObject result;
    for (String fpath : fileList) {
      success = true;
      error = "";
      try {
        parser.parse(fpath);
      } catch (Exception e) {
        success = false;
        error = e.getMessage();
      }
      result = new JSONObject();
      result.put("file", fpath.replaceAll(exDir, ""));
      result.put("success", success);
      result.put("error", error);
      results.add(result);
    }

    report.put("results", results);
    Utils.saveReport(report, outdir);
  }

  public static void main(String[] args) {
    CommandLine.run(new RamlTckRunner(), args);
  }
}
