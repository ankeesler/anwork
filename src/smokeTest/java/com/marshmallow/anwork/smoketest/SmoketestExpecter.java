package com.marshmallow.anwork.smoketest;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.Reader;
import java.lang.ProcessBuilder.Redirect;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;
import java.util.stream.Collectors;
import java.util.stream.Stream;

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
   * Match a multiple regular expressions against the output of a CLI command, or
   * {@link Assert#fail()}. Each line of CLI output will be individually matched against the next
   * {@code expectRegexes}. If all {@code expectRegexes} do not find a matching line, then this
   * method will {@link Assert#fail()}.
   *
   * @param expectedRegexes The regular expressions to try to find in the CLI command output
   * @param processBuilder The {@link ProcessBuilder} that contains the necessary commands (see
   *     {@link ProcessBuilder#command()}) and working directory (see
   *     {@link ProcessBuilder#directory()}) to be run
   * @throws Exception if something goes wrong
   */
  public static void expect(String[] expectedRegexes, ProcessBuilder processBuilder)
      throws Exception {
    processBuilder.redirectOutput(Redirect.PIPE);

    Process process = processBuilder.start();
    int exitValue = runProcess(process);
    assertEquals("Process " + processBuilder.command() + " failed with exit code " + exitValue,
                 0, exitValue);

    List<String> cliOutputLines = getCliOutputLines(process);
    Map<String, String> matches = findFirstMatches(expectedRegexes, cliOutputLines);
    if (matches == null) {
      String message = ("Could not find regexes " + Arrays.toString(expectedRegexes) + " "
                        + "in output lines:\n" + makePrettyLines(cliOutputLines));
      System.out.println(message);
      fail(message);
    } else {
      for (String regex : matches.keySet()) {
        System.out.println("Matched '" + regex + "' to '" + matches.get(regex));
      }
    }
  }

  /**
   * Assert that a number of regular expressions do not all match some CLI output or else,
   * {@link Assert#fail()}. Each line of CLI output will be individually matched against the next
   * {@code expectRegexes}. If all {@code expectRegexes} do find a matching line, then this
   * method will {@link Assert#fail()}.
   *
   * @param nexpectedRegexes The regular expressions to hope to not find in the CLI command output
   * @param processBuilder The {@link ProcessBuilder} that contains the necessary commands (see
   *     {@link ProcessBuilder#command()}) and working directory (see
   *     {@link ProcessBuilder#directory()}) to be run
   * @throws Exception if something goes wrong
   */
  public static void nexpect(String[] nexpectedRegexes, ProcessBuilder processBuilder)
      throws Exception {
    processBuilder.redirectOutput(Redirect.PIPE);

    Process process = processBuilder.start();
    int exitValue = runProcess(process);
    assertEquals("Process " + processBuilder.command() + " failed with exit code " + exitValue,
                 0, exitValue);

    List<String> cliOutputLines = getCliOutputLines(process);
    Map<String, String> matches = findFirstMatches(nexpectedRegexes, cliOutputLines);
    if (matches != null) {
      String message = ("Found regexes " + Arrays.toString(nexpectedRegexes) + " "
                        + "in output lines:\n" + makePrettyLines(cliOutputLines));
      System.out.println(message);
      fail(message);
    } else {
      for (String regex : nexpectedRegexes) {
        System.out.println("Did not get " + regex + " in " + cliOutputLines);
      }
    }
  }

  // Returns a map from the regex to the line it matched, or null if the match failed
  private static Map<String, String> findFirstMatches(String[] expectRegexes, List<String> lines) {
    if (expectRegexes == null || expectRegexes.length == 0 || expectRegexes[0] == null) {
      return Collections.emptyMap();
    } else if (lines.size() == 0) {
      return null;
    }

    Map<String, String> matches = new LinkedHashMap<String, String>();
    Pattern[] patterns = Stream.of(expectRegexes).map(Pattern::compile).toArray(Pattern[]::new);
    int lineIndex = 0;
    int patternIndex = 0;
    do {
      String line = lines.get(lineIndex);
      Pattern pattern = patterns[patternIndex];

      if (pattern.matcher(line).matches()) {
        matches.put(pattern.pattern(), line);
        patternIndex++;
      }

      lineIndex++;
    } while (lineIndex < lines.size() && patternIndex < patterns.length);

    return (patternIndex == patterns.length ? matches : null);
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

  private static String makePrettyLines(List<String> lines) {
    return (lines.size() == 0
           ? "<no output>"
           : lines.stream().collect(Collectors.joining("\n")));
  }
}
