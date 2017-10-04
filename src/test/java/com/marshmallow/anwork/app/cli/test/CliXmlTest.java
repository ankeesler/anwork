package com.marshmallow.anwork.app.cli.test;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;
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

  private Cli cli;

  @Test
  public void testGood() throws Exception {
    cli = read("cli-xml-test.xml");
    assertNotNull(cli);

    // Reset the run count on the test action.
    BringHomeBaconTestCliAction.resetRunCount();

    // Test flags.
    run("-a", "--bacon", "-c", "fish", "--dog", "97");
    run("-a", "-b", "-c", "fish", "-d", "34");
    assertEquals(0, BringHomeBaconTestCliAction.getRunCount());

    // Test commands.
    run("fish");
    run("marlin");
    assertEquals(0, BringHomeBaconTestCliAction.getRunCount());

    // Test commands and flags.
    run("-a", "--bacon", "-c", "fish", "--dog", "25", "fish");
    run("-a", "-b", "-c", "fish", "-d", "5", "fish");
    run("-a", "--bacon", "-c", "fish", "--dog", "15", "marlin");
    run("-a", "-b", "-c", "fish", "-d", "35", "marlin");
    assertEquals(0, BringHomeBaconTestCliAction.getRunCount());

    // Test lists.
    run("list-a", "-m", "--dad", "moving-the-grass", "bring-home-bacon");
    run("list-b", "shake-it-up");
    assertEquals(1, BringHomeBaconTestCliAction.getRunCount());
    assertArrayEquals(new String[0], BringHomeBaconTestCliAction.getRunArguments(0));

    // Test list commands with arguments.
    run("list-a", "-m", "--dad", "moving-the-grass", "bring-home-bacon", "hey", "ho");
    assertEquals(2, BringHomeBaconTestCliAction.getRunCount());
    assertArrayEquals(new String[] { "hey", "ho" },
                      BringHomeBaconTestCliAction.getRunArguments(1));

    // Test lists and flags.
    run("-a", "--bacon", "-c", "fish", "--dog", "1", "list-a");
    run("-a", "-b", "-c", "fish", "-d", "2", "list-a");
    run("-a", "--bacon", "-c", "fish", "--dog", "3", "list-b");
    run("-a", "-b", "-c", "fish", "-d", "4", "list-b");
    assertEquals(2, BringHomeBaconTestCliAction.getRunCount());
  }

  @Test(expected = Exception.class)
  public void testMissingCliTag() throws Exception {
    read("missing-cli-tag.xml");
  }

  @Test(expected = Exception.class)
  public void testUnknownParameterType() throws Exception {
    read("bad-parameter-type.xml");
  }

  @Test(expected = Exception.class)
  public void testBadClassName() throws Exception {
    read("bad-class-name.xml");
  }

  @Test(expected = Exception.class)
  public void testUnknownClass() throws Exception {
    read("unknown-class.xml");
  }

  @Test(expected = Exception.class)
  public void testBadClassType() throws Exception {
    read("bad-class-type.xml");
  }

  private Cli read(String filename) throws Exception {
    File file = TestUtilities.getFile(filename, getClass());
    return new CliXmlReader(file).read();
  }

  private void run(String...args) {
    cli.parse(args);
  }
}
