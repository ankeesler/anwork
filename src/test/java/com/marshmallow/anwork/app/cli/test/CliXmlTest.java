package com.marshmallow.anwork.app.cli.test;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliXmlReader;
import com.marshmallow.anwork.core.test.TestUtilities;

import java.io.File;
import java.io.FileInputStream;
import java.io.InputStream;

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
public class CliXmlTest extends BaseCliTest {

  @Override
  protected Cli createCli() throws Exception {
    BringHomeBaconTestCliAction.resetRunCount();
    TestCliActionCreator.resetCreatedActions();
    return read("cli-xml-test.xml");
  }

  /*
   * Section - Positive Tests
   */

  @Test
  public void testFlags() {
    parse("-a", "--bacon", "-c", "fish", "--dog", "97");
    parse("-a", "-b", "-c", "fish", "-d", "34");
    assertEquals(0, BringHomeBaconTestCliAction.getRunCount());
    TestCliAction fishAction = TestCliActionCreator.getCreatedAction("fish");
    assertNotNull(fishAction);
    assertFalse(fishAction.getRan());
    TestCliAction marlinAction = TestCliActionCreator.getCreatedAction("marlin");
    assertNotNull(marlinAction);
    assertFalse(marlinAction.getRan());
  }

  @Test
  public void testCommands() {
    parse("fish");
    parse("marlin");
    assertEquals(0, BringHomeBaconTestCliAction.getRunCount());
    TestCliAction fishAction = TestCliActionCreator.getCreatedAction("fish");
    assertNotNull(fishAction);
    assertTrue(fishAction.getRan());
    TestCliAction marlinAction = TestCliActionCreator.getCreatedAction("marlin");
    assertNotNull(marlinAction);
    assertTrue(marlinAction.getRan());
  }

  @Test
  public void testCommandsAndFlags() {
    parse("-a", "--bacon", "-c", "fish", "--dog", "25", "fish");
    parse("-a", "-b", "-c", "fish", "-d", "5", "fish");
    parse("-a", "--bacon", "-c", "fish", "--dog", "15", "marlin");
    parse("-a", "-b", "-c", "fish", "-d", "35", "marlin");
    assertEquals(0, BringHomeBaconTestCliAction.getRunCount());
  }

  @Test
  public void testLists() {
    parse("list-a", "-m", "--dad", "moving-the-grass", "bring-home-bacon");
    parse("list-b", "shake-it-up");
    assertEquals(1, BringHomeBaconTestCliAction.getRunCount());
    assertArrayEquals(new String[0], BringHomeBaconTestCliAction.getRunArguments(0));
  }

  @Test
  public void testListCommandsWithArguments() {
    parse("list-a", "-m", "--dad", "moving-the-grass", "bring-home-bacon", "hey", "ho");
    assertEquals(1, BringHomeBaconTestCliAction.getRunCount());
    assertArrayEquals(new String[] { "hey", "ho" },
                      BringHomeBaconTestCliAction.getRunArguments(0));
  }

  @Test
  public void testListsAndFlags() {
    parse("-a", "--bacon", "-c", "fish", "--dog", "1", "list-a");
    parse("-a", "-b", "-c", "fish", "-d", "2", "list-a");
    parse("-a", "--bacon", "-c", "fish", "--dog", "3", "list-b");
    parse("-a", "-b", "-c", "fish", "-d", "4", "list-b");
    assertEquals(0, BringHomeBaconTestCliAction.getRunCount());
  }

  /*
   * Section - Negative Tests
   */

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

  @Test(expected = Exception.class)
  public void testNoActionOrActionCreator() throws Exception {
    read("no-action-or-action-creator.xml");
  }

  @Test(expected = Exception.class)
  public void testBothActionAndActionCreator() throws Exception {
    read("both-action-and-action-creator.xml");
  }

  private Cli read(String filename) throws Exception {
    File file = TestUtilities.getFile(filename, getClass());
    InputStream stream = new FileInputStream(file);
    return new CliXmlReader(stream).read();
  }
}
