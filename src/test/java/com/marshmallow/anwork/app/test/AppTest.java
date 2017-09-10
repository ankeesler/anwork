package com.marshmallow.anwork.app.test;

import com.marshmallow.anwork.app.AnworkApp;
import com.marshmallow.anwork.core.test.TestUtilities;

import java.io.File;

import org.junit.Test;

/**
 * This is a test for the anwork application.
 *
 * @author Andrew
 * Created Sep 9, 2017
 */
public class AppTest {

  private static final String CONTEXT = "app-test-context";
  private static final String PERSISTENCE_ROOT
    = new File(TestUtilities.TEST_RESOURCES_ROOT, "app-test").getAbsolutePath();

  @Test
  public void runTest() {
    runApp(new String[] { "-d",
                          "--context", CONTEXT,
                          "-o", PERSISTENCE_ROOT,
                          "-s" });
  }

  private void runApp(String[] args) {
    AnworkApp.main(args);
  }
}
