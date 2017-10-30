package com.marshmallow.anwork.app.cli.test;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.app.cli.ArgumentType;
import com.marshmallow.anwork.app.cli.ArgumentValues;
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
    parse("-a", "--bacon", "-c", "fish", "--dog", "97", "--no-description-long-flag", "meh");
    parse("-a", "-b", "-c", "fish", "-d", "34", "-no-description-short-flag", "-n", "meh");
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
    parse("fish", "25", "steve");
    parse("fish", "-o", "15", "25", "steve");
    parse("marlin");
    assertEquals(0, BringHomeBaconTestCliAction.getRunCount());
    TestCliAction fishAction = TestCliActionCreator.getCreatedAction("fish");
    assertNotNull(fishAction);
    assertTrue(fishAction.getRan());
    ArgumentValues arguments = fishAction.getArguments();
    assertEquals(2, arguments.getAllKeys().length);
    assertTrue(arguments.containsKey("number"));
    assertTrue(arguments.containsKey("name"));
    assertEquals(new Long(25), arguments.getValue("number", ArgumentType.NUMBER));
    assertEquals("steve", arguments.getValue("name", ArgumentType.STRING));
    assertNull(arguments.getValue("wrong", ArgumentType.STRING));
    TestCliAction marlinAction = TestCliActionCreator.getCreatedAction("marlin");
    assertNotNull(marlinAction);
    assertTrue(marlinAction.getRan());
  }

  @Test
  public void testCommandsAndFlags() {
    parse("-a", "--bacon", "-c", "fish", "--dog", "25", "fish", "25", "steve");
    parse("-a", "-b", "-c", "fish", "-d", "5", "fish", "25", "steve");
    parse("-a", "--bacon", "-c", "fish", "--dog", "15", "marlin");
    parse("-a", "-b", "-c", "fish", "-d", "35", "marlin");
    assertEquals(0, BringHomeBaconTestCliAction.getRunCount());
  }

  @Test
  public void testLists() {
    parse("list-a", "-m", "--dad", "moving-the-grass", "bring-home-bacon", "hey", "ho");
    parse("list-b", "shake-it-up", "--andrew", "-o", "foo");
    assertEquals(1, BringHomeBaconTestCliAction.getRunCount());
    ArgumentValues arguments = BringHomeBaconTestCliAction.getRunArguments(0);
    assertEquals(2, arguments.getAllKeys().length);
    assertEquals("hey", arguments.getValue("arg0", ArgumentType.STRING));
    assertEquals("ho", arguments.getValue("arg1", ArgumentType.STRING));
  }

  @Test
  public void testListCommandsWithArguments() {
    parse("list-a", "-m", "--dad", "moving-the-grass", "bring-home-bacon", "hey", "ho");
    assertEquals(1, BringHomeBaconTestCliAction.getRunCount());
    ArgumentValues arguments = BringHomeBaconTestCliAction.getRunArguments(0);
    assertEquals(2, arguments.getAllKeys().length);
    assertEquals("hey", arguments.getValue("arg0", ArgumentType.STRING));
    assertEquals("ho", arguments.getValue("arg1", ArgumentType.STRING));
  }

  @Test
  public void testListsAndFlags() {
    parse("-a", "--bacon", "-c", "fish", "--dog", "1", "list-a");
    parse("-a", "-b", "-c", "fish", "-d", "2", "list-a");
    parse("-a", "--bacon", "-c", "fish", "--dog", "3", "list-b");
    parse("-a", "-b", "-c", "fish", "-d", "4", "list-b");
    assertEquals(0, BringHomeBaconTestCliAction.getRunCount());
  }

  @Test
  public void testVisitorChecksOut() {
    TestCliVisitor visitor = new TestCliVisitor();
    visit(visitor);
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedShortFlags(),
                                            "a", "no-description-short-flag");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedShortFlagsWithParameters(),
                                            "c", "e", "o");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedLongFlags(),
                                            "bacon", "mom", "andrew");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedLongFlagsWithParameters(),
                                            "dog", "no-description-long-flag", "dad", "output");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedCommands(),
                                            "fish", "marlin", "bring-home-bacon", "shake-it-up");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedCommandArguments(),
                                            "number", "name", "arg0", "arg1");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedLists(),
                                            "tuna", "list-a", "dumb-list", "list-b");
    TestUtilities.assertVariadicArrayEquals(visitor.getLeftLists(),
                                            "dumb-list", "list-a", "list-b", "tuna");
  }

  @Test
  public void testOptionalData() {
    OptionalDataCliVisitor visitor = new OptionalDataCliVisitor();
    visit(visitor);
    TestUtilities.assertVariadicArrayEquals(visitor.getFlagsWithDescriptions(),
                                            "a", "b", "c", "d", "e", "o", "m", "o");
    TestUtilities.assertVariadicArrayEquals(visitor.getCommandsWithDescriptions(),
                                            "fish", "marlin", "shake-it-up");
    TestUtilities.assertVariadicArrayEquals(visitor.getCommandArgumentsWithDescriptions(),
                                            "number");
    TestUtilities.assertVariadicArrayEquals(visitor.getListsWithDescriptions(),
                                            "tuna", "list-a", "list-b");
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

  @Test(expected = Exception.class)
  public void testTooManyFlagArguments() throws Exception {
    read("too-many-flag-arguments.xml");
  }

  private Cli read(String filename) throws Exception {
    File file = TestUtilities.getFile(filename, getClass());
    InputStream stream = new FileInputStream(file);
    return new CliXmlReader(stream).read();
  }
}
