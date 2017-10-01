package com.marshmallow.anwork.app.cli.test;

import static org.junit.Assert.assertNotNull;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliXmlReader;
import com.marshmallow.anwork.core.test.TestUtilities;

import java.io.File;

import org.junit.Test;

/**
 * This is a test for stuff related to CLI XML files.
 *
 * <p>
 * Created Oct 1, 2017
 * </p>
 *
 * @author Andrew
 */
public class CliXmlTest {

  private final static File TEST_RESOURCE_ROOT
      = new File(TestUtilities.TEST_RESOURCES_ROOT, "cli-xml-test");

  private Cli cli;

  @Test
  public void testGood() throws Exception {
    File file = new File(TEST_RESOURCE_ROOT, "cli-xml-test.xml");
    CliXmlReader reader = new CliXmlReader(file);
    cli = reader.read();
    assertNotNull(cli);

    // Test flags.
    run("-a", "--bacon", "-c", "fish", "--dog", "marlin");
    run("-a", "-b", "-c", "fish", "-d", "marlin");

    // Test commands.
    run("fish");
    run("marlin");

    // Test commands and flags.
    run("-a", "--bacon", "-c", "fish", "--dog", "marlin", "fish");
    run("-a", "-b", "-c", "fish", "-d", "marlin", "fish");
    run("-a", "--bacon", "-c", "fish", "--dog", "marlin", "marlin");
    run("-a", "-b", "-c", "fish", "-d", "marlin", "marlin");

    // Test lists.
    run("list-a", "-m", "--dad", "moving-the-grass", "bring-home-bacon");
    run("list-b", "shake-it-up");

    // Test lists and flags.
    run("-a", "--bacon", "-c", "fish", "--dog", "marlin", "list-a");
    run("-a", "-b", "-c", "fish", "-d", "marlin", "list-a");
    run("-a", "--bacon", "-c", "fish", "--dog", "marlin", "list-b");
    run("-a", "-b", "-c", "fish", "-d", "marlin", "list-b");
  }

  @Test(expected = Exception.class)
  public void testBad() throws Exception {
    File file = new File(TEST_RESOURCE_ROOT, "missing-cli-tag.xml");
    new CliXmlReader(file).read();
  }

  private void run(String...args) {
    cli.parse(args);
  }
}
