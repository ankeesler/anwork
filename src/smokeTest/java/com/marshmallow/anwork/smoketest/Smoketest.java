package com.marshmallow.anwork.smoketest;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.fail;

import java.io.File;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

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

  @Test
  public void createTest() throws Exception {
    run(true, "create", "task-a");
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
