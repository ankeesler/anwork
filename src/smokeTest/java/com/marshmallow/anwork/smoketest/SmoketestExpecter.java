package com.marshmallow.anwork.smoketest;

import java.io.BufferedReader;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.Reader;
import java.lang.ProcessBuilder.Redirect;
import java.util.ArrayList;
import java.util.List;
import java.util.regex.Pattern;
import java.util.stream.Collectors;
import java.util.stream.Stream;

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

  private static final boolean DEBUG = true;

  // Use #expect below.
  private SmoketestExpecter() { }

  /**
   * Match a multiple regular expressions against the output of a CLI command. Each line of CLI
   * output will be individually matched against the next {@code expectRegexes}. This method will
   * return an ordered list of lines of output that were matches against the {@code expectRegexes}.
   *
   * <p>Note that this means that if the length of the returned array does not match the length of the
   * passed {@code expectRegexes}, then at least one provided {@code expectRegexes} was not found
   * in the command output. Here is a sample call to this function.
   * <pre>
   * ProcessBuilder processBuilder = makeProcessBuilder();
   * String[] expected = new String[] { ".*foo.*", "^bar$" };
   * String[] found = SmoketestExpecter.expect(expected, processBuilder);
   * if (found.length != expected.length) {
   *   // fail!
   *   String didntFind = expected[found.length];
   *   ...
   * } else {
   *   // success!
   *   ...
   * }
   * </pre>
   *
   * @param expectedRegexes The regular expressions to try to find in the CLI command output
   * @param processBuilder The {@link ProcessBuilder} that contains the necessary commands (see
   *     {@link ProcessBuilder#command()}) and working directory (see
   *     {@link ProcessBuilder#directory()}) to be run
   * @return An ordered list of lines of output that were matches against the {@code expectRegexes}
   * @throws Exception if something goes wrong
   */
  public static String[] expect(String[] expectedRegexes, ProcessBuilder processBuilder)
      throws Exception {
    processBuilder.redirectOutput(Redirect.PIPE);

    Process process = processBuilder.start();
    int exitValue = runProcess(process);
    if (exitValue != 0) {
      throw new Exception("Process " + processBuilder.command()
                          + " failed with exit code " + exitValue);
    }

    List<String> cliOutputLines = getCliOutputLines(process);
    String[] matches = findFirstMatches(expectedRegexes, cliOutputLines);
    if (DEBUG) {
      for (int i = 0; i < expectedRegexes.length; i++) {
        if (i < matches.length) {
          System.out.println("Matched '" + expectedRegexes[i] + "' to '" + matches[i]);
        } else {
          System.out.println("Did not find regex '" + "' in output lines:\n"
                             + makePrettyLines(cliOutputLines));
        }
      }
    }
    return matches;
  }

  // Returns a map from the regex to the line it matched, or null if the match failed
  private static String[] findFirstMatches(String[] expectRegexes, List<String> lines) {
    if (expectRegexes == null
        || expectRegexes.length == 0
        || expectRegexes[0] == null
        || lines.size() == 0) {
      return new String[0];
    }

    List<String> matches = new ArrayList<String>(expectRegexes.length);
    Pattern[] patterns = Stream.of(expectRegexes).map(Pattern::compile).toArray(Pattern[]::new);
    int lineIndex = 0;
    int patternIndex = 0;
    do {
      String line = lines.get(lineIndex);
      Pattern pattern = patterns[patternIndex];

      if (pattern.matcher(line).matches()) {
        matches.add(line);
        patternIndex++;
      }

      lineIndex++;
    } while (lineIndex < lines.size() && patternIndex < patterns.length);

    return matches.toArray(new String[0]);
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
        throw new Exception("Process " + process + " took longer than 1 second"
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
