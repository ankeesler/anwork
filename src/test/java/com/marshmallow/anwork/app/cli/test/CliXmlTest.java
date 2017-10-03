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

  private static final File TEST_RESOURCE_ROOT
      = new File(TestUtilities.TEST_RESOURCES_ROOT, "cli-xml-test");

  private Cli cli;

  @Test
  public void testGood() throws Exception {
    File file = new File(TEST_RESOURCE_ROOT, "cli-xml-test.xml");
    CliXmlReader reader = new CliXmlReader(file);
    cli = reader.read();
    assertNotNull(cli);

    // Test flags.
    run("-a", "--bacon", "-c", "fish", "--dog", "97");
    run("-a", "-b", "-c", "fish", "-d", "34");

    // Test commands.
    run("fish");
    run("marlin");

    // Test commands and flags.
    run("-a", "--bacon", "-c", "fish", "--dog", "25", "fish");
    run("-a", "-b", "-c", "fish", "-d", "5", "fish");
    run("-a", "--bacon", "-c", "fish", "--dog", "15", "marlin");
    run("-a", "-b", "-c", "fish", "-d", "35", "marlin");

    // Test lists.
    run("list-a", "-m", "--dad", "moving-the-grass", "bring-home-bacon");
    run("list-b", "shake-it-up");

    // Test lists and flags.
    run("-a", "--bacon", "-c", "fish", "--dog", "1", "list-a");
    run("-a", "-b", "-c", "fish", "-d", "2", "list-a");
    run("-a", "--bacon", "-c", "fish", "--dog", "3", "list-b");
    run("-a", "-b", "-c", "fish", "-d", "4", "list-b");
  }

  @Test(expected = Exception.class)
  public void testMissingCliTag() throws Exception {
    File file = new File(TEST_RESOURCE_ROOT, "missing-cli-tag.xml");
    new CliXmlReader(file).read();
  }

  @Test(expected = Exception.class)
  public void testUnknownParameterType() throws Exception {
    File file = new File(TEST_RESOURCE_ROOT, "bad-parameter-type.xml");
    new CliXmlReader(file).read();
  }

  private void run(String...args) {
    cli.parse(args);
  }
}
