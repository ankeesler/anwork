package com.marshmallow.anwork.smoketest;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.fail;

import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.Reader;
import java.lang.ProcessBuilder.Redirect;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

import org.junit.Assert;

/**
 * This is a class that runs some ANWORK CLI commands and expects some regular expression of
 * output.
 *
 * <p>
 * Created Oct 30, 2017
 * </p>
 *
 * @author Andrew
 */
public class SmoketestExpecter {

  // Use #expect below.
  private SmoketestExpecter() { }

  /**
   * Match a regular expression against the output of a CLI command, or {@link Assert#fail()}. Each
   * line of CLI output will be individually matched against the {@code expectRegex}.
   *
   * @param expectRegex The regular expression to try to find in the CLI command output
   * @param processBuilder The {@link ProcessBuilder} that contains the necessary commands (see
   *     {@link ProcessBuilder#command()}) and working directory (see
   *     {@link ProcessBuilder#directory()}) to be run
   * @throws Exception if something goes wrong
   */
  public static void expect(String expectRegex, ProcessBuilder processBuilder) throws Exception {
    processBuilder.redirectOutput(Redirect.PIPE);

    Process process = processBuilder.start();
    int exitValue = runProcess(process);
    assertEquals("Process " + processBuilder.command() + " failed with exit code " + exitValue,
                 0, exitValue);

    List<String> cliOutputLines = getCliOutputLines(process);
    String firstMatch = findFirstMatch(expectRegex, cliOutputLines);
    assertNotNull(("Could not find regex '" + expectRegex + "' "
                   + "in output lines " + cliOutputLines),
                  firstMatch);
  }

  private static String findFirstMatch(String expectRegex, List<String> lines) {
    if (expectRegex == null) {
      return "";
    }

    Pattern pattern = Pattern.compile(expectRegex);
    for (String line : lines) {
      Matcher matcher = pattern.matcher(line);
      if (matcher.matches()) {
        return line;
      }
    }
    return null;
  }

  // Runs the process for a maximum of 1 second before saying that it failed
  // Returns the exit code from the underlying process
  private static int runProcess(Process process) throws Exception {
    long startMillis = System.currentTimeMillis();
    while (process.isAlive()) {
      long nowMillis = System.currentTimeMillis();
      if (nowMillis - startMillis > 1000) {
        process.destroy();
        if (process.isAlive()) {
          process.destroyForcibly();
        }
        fail("Process " + process + " took longer than 1 second"
             + " and was destroyed.");
      }
    }
    return process.exitValue();
  }

  private static List<String> getCliOutputLines(Process process) throws Exception {
    InputStream processOutput = process.getInputStream();
    List<String> lines = new ArrayList<String>();
    try (Reader reader = new InputStreamReader(processOutput)) {
      try (BufferedReader lineReader = new BufferedReader(reader)) {
        String line;
        while ((line = lineReader.readLine()) != null) {
          lines.add(line);
        }
      }
    }
    return lines;
  }
}
