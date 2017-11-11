package com.marshmallow.anwork.smoketest;

import static org.junit.Assert.fail;
import static org.junit.Assert.assertEquals;

import java.io.File;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import org.junit.Before;
import org.junit.BeforeClass;
import org.junit.Test;

/**
 * This is a system level test fixture that runs using packaged ANWORK app.
 *
 * <p>
 * Created Oct 5, 2017
 * </p>
 *
 * @author Andrew
 */
public class Smoketest {

  private static final String ANWORK_SMOKETEST_DIR_ENV_VAR = "ANWORK_SMOKETEST_DIR";

  private static File smoketestDirectory = null;
  private static File anworkBinary = null;

  /**
   * Initialize the static fields in this class that indicate the ANWORK package to test.
   */
  @BeforeClass
  public static void findAnworkPackage() {
    String smoketestDir = System.getenv(ANWORK_SMOKETEST_DIR_ENV_VAR);
    if (smoketestDir == null) {
      fail("Cannot find smoketest directory from environmental variable "
           + ANWORK_SMOKETEST_DIR_ENV_VAR);
    }

    smoketestDirectory = new File(smoketestDir);
    if (!smoketestDirectory.exists()) {
      fail("Directory " + smoketestDirectory + " for ANWORK smoketest does not exist!");
    }

    anworkBinary = new File(smoketestDirectory, "anwork/bin/anwork");
    if (!anworkBinary.exists()) {
      fail("Binary " + anworkBinary + " from ANWORK package does not exist!");
    }
  }

  @Before
  public void deleteContext() throws Exception {
    ProcessBuilder processBuilder = new ProcessBuilder();
    processBuilder.directory(smoketestDirectory);
    processBuilder.command("rm", "-f", "default-context");
    Process process = processBuilder.start();
    while (process.isAlive()) {
      // wait...
    }
    assertEquals(0, process.exitValue());
  }

  @Test
  public void showNoTasksTest() throws Exception {
    run("task", "show");
  }

  @Test
  public void showSomeTasksTest() throws Exception {
    run("task", "create", "task-a",
        "-e", "This is task-a",
        "-p", "1");
    run("task", "create", "task-b",
        "--description", "This is task-b",
        "--priority", "1");
    run("task", "create", "task-c");
    expect(new String[] { "task", "show" },
           new String[] { "RUNNING tasks:",
                          "BLOCKED tasks:",
                          "WAITING tasks:",
                          "  task-a.*",
                          "  task-b .*",
                          "  task-c.*",
                          "FINISHED tasks:", });
    nexpect(new String[] { "task", "show" },
            new String[] { "  task-b \\(1\\).*",
                           "  task-c \\(1\\).*"});
    nexpect(new String[] { "task", "show", "--short", },
            new String[] { "    created.*",
                           "    priority.*", });
  }

  @Test
  public void deleteSomeTasksTest() throws Exception {
    run("task", "create", "task-a");
    run("task", "create", "task-b");
    run("task", "delete", "task-a");
    expect(new String[] { "task", "show" },
           new String[] { "WAITING tasks:", "  task-b.*"});
  }

  @Test
  public void setStateOnSomeTasks() throws Exception {
    run("task", "create", "task-a");
    run("task", "create", "task-b");
    run("task", "create", "task-c");
    run("task", "set-running", "task-c");
    run("task", "set-blocked", "task-b");
    run("task", "set-finished", "task-a");
    expect(new String[] { "task", "show" },
           new String[] { "RUNNING tasks:", "  task-c.*",
                          "BLOCKED tasks:", "  task-b.*",
                          "FINISHED tasks:", "  task-a.*"});
  }

  @Test
  public void showAllEmptyJournalTest() throws Exception {
    run("journal", "show-all");
  }

  @Test
  public void showJournalTest() throws Exception {
    run("task", "create", "task-a");
    run("task", "create", "task-b");
    run("task", "create", "task-c");
    run("journal", "show", "task-a");
    run("journal", "show", "task-b");
    run("journal", "show", "task-c");
  }

  @Test
  public void setPriorityTest() throws Exception {
    run("task", "create", "task-a");
    run("task", "create", "task-b");
    run("task", "set-priority", "task-a", "20");
    run("task", "set-priority", "task-b", "21");
    run("task", "set-priority", "task-a", "22");
    expect(new String[] { "task", "show", },
           new String[] { "WAITING tasks:",
                          "  task-b \\(\\d+\\)",
                          "    priority 21",
                          "  task-a \\(\\d+\\)",
                          "    priority 22", });
  }

  @Test
  public void summaryTest() throws Exception {
    run("task", "create", "task-a");
    run("task", "create", "task-b");
    run("task", "set-running", "task-b");
    run("task", "set-blocked", "task-b");
    run("task", "note", "task-b", "tuna");
    run("task", "set-running", "task-a");
    run("task", "set-finished", "task-a");
    run("task", "set-running", "task-b");
    run("task", "note", "task-b", "fish");
    run("task", "set-finished", "task-b");
    expect(new String[] { "summary", "2" },
           new String[] { "\\[.*\\]: Finished 'task-b'",
                          "  took \\d+ seconds",
                          "\\[.*\\]: Finished 'task-a'",
                          "  took \\d+ seconds", });
  }

  @Test
  public void makeSureDebugPrintingWorks() throws Exception {
    expect(new String[] { "-d", "task", "create", "task-a" },
           new String[] { ".*created task.*"});
  }

  // We need both this method and makeSureTasksAreActuallyDeletedBetweenMethodsPart2 so that we can
  // ensure we are running this check _after_ the first @Test method has run.
  @Test
  public void makeSureTasksAreActuallyDeletedBetweenMethodsPart1() throws Exception {
    nexpect(new String[] { "journal", "show-all" },
            new String[] { ".*", } );
    run("task", "create", "task-a");
  }

  // We need both this method and makeSureTasksAreActuallyDeletedBetweenMethodsPart1 so that we can
  // ensure we are running this check _after_ the first @Test method has run.
  @Test
  public void makeSureTasksAreActuallyDeletedBetweenMethodsPart2() throws Exception {
    nexpect(new String[] { "journal", "show-all" },
            new String[] { ".*", } );
    run("task", "create", "task-a");
  }

  private void run(String...args) throws Exception {
    expect(args, new String[] { null });
  }

  private void expect(String[] args, String[] expectedRegexes) throws Exception {
    expectOrNexpect(args, expectedRegexes, true);
  }

  private void nexpect(String[] args, String[] nexpectedRegexes) throws Exception {
    expectOrNexpect(args, nexpectedRegexes, false);
  }

  private void expectOrNexpect(String[] args, String[] expectRegexes, boolean expect)
      throws Exception {
    List<String> commands = new ArrayList<String>();
    commands.add(anworkBinary.getAbsolutePath());
    commands.addAll(Arrays.asList(args));

    ProcessBuilder processBuilder = new ProcessBuilder(commands);
    configureProcess(processBuilder);
    if (expect) {
      SmoketestExpecter.expect(expectRegexes, processBuilder);
    } else {
      SmoketestExpecter.nexpect(expectRegexes, processBuilder);
    }
  }

  private void configureProcess(ProcessBuilder processBuilder) {
    processBuilder.directory(smoketestDirectory);
    processBuilder.inheritIO();
  }
}
