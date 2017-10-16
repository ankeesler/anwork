package com.marshmallow.anwork.smoketest;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

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
  public void deleteAllTasks() throws Exception {
    run(true, "task", "delete-all");
  }

  @Test
  public void showNoTasksTest() throws Exception {
    run(true, "task", "show");
  }

  @Test
  public void showSomeTasksTest() throws Exception {
    run(true, "task", "create", "task-a", "This is task-a", "1");
    run(true, "task", "create", "task-b", "This is task-b", "1");
    run(true, "task", "create", "task-c", "This is task-c", "1");
    run(true, "task", "show");
  }

  @Test
  public void deleteSomeTasksTest() throws Exception {
    run(true, "task", "create", "task-a", "This is task-a", "1");
    run(true, "task", "create", "task-b", "This is task-b", "1");
    run(true, "task", "delete", "task-a");
    run(true, "task", "show");
  }

  @Test
  public void setStateOnSomeTasks() throws Exception {
    run(true, "task", "create", "task-a", "This is task-a", "1");
    run(true, "task", "create", "task-b", "This is task-b", "1");
    run(true, "task", "create", "task-c", "This is task-c", "1");
    run(true, "task", "set-running", "task-c");
    run(true, "task", "set-blocked", "task-b");
    run(true, "task", "set-finished", "task-a");
    run(true, "task", "show");
  }

  @Test
  public void showAllEmptyJournalTest() throws Exception {
    run(true, "journal", "show-all");
  }

  @Test
  public void showJournalTest() throws Exception {
    run(true, "task", "create", "task-a", "This is task-a", "1");
    run(true, "task", "create", "task-b", "This is task-b", "1");
    run(true, "task", "create", "task-c", "This is task-c", "1");
    run(true, "journal", "show", "task-a");
    run(true, "journal", "show", "task-b");
    run(true, "journal", "show", "task-c");
  }

  private void run(boolean debug, String...args) throws Exception {
    List<String> commands = new ArrayList<String>();
    commands.add(anworkBinary.getAbsolutePath());
    if (debug) {
      commands.add("-d");
    }
    commands.addAll(Arrays.asList(args));

    ProcessBuilder processBuilder = new ProcessBuilder(commands);
    configureProcess(processBuilder);
    int exitCode = runProcess(processBuilder);
    assertEquals("ANWORK process " + processBuilder.command()
                 + " failed with exit code " + exitCode, 0, exitCode);
  }

  private void configureProcess(ProcessBuilder processBuilder) {
    processBuilder.directory(smoketestDirectory);
    processBuilder.inheritIO();
  }

  // Runs the process for a maximum of 1 second before saying that it failed
  // Returns the exit code from the underlying process
  private int runProcess(ProcessBuilder processBuilder) throws Exception {
    System.out.println("Running " + processBuilder.command() + "...");
    Process process = processBuilder.start();
    long startMillis = System.currentTimeMillis();
    while (process.isAlive()) {
      long nowMillis = System.currentTimeMillis();
      if (nowMillis - startMillis > 1000) {
        process.destroy();
        if (process.isAlive()) {
          process.destroyForcibly();
        }
        fail("Process " + processBuilder.command() + " took longer than 1 second"
             + " and was destroyed.");
      }
    }
    return process.exitValue();
  }
}
