package com.marshmallow.anwork.app.cli.test;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliArgumentType;
import com.marshmallow.anwork.app.cli.CliCommand;
import com.marshmallow.anwork.app.cli.CliFlags;
import com.marshmallow.anwork.app.cli.CliList;

import java.util.List;

import org.junit.Test;

/**
 * This is a {@link BaseCliTest} that tests basic CLI functionality on a {@link Cli} created via
 * Java code.
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public class BasicCliTest extends BaseCliTest {

  private TestCliAction tunaMarlinAction = new TestCliAction();

  private TestCliAction mayoAction = new TestCliAction();

  @Override
  protected Cli createCli() {
    Cli cli = new Cli("cli-test",  "The are commands for the CLI unit test");
    CliList root = cli.getRoot();
    root.addShortFlag("a",
                      "Description for a flag");
    root.addLongFlag("b",
                     "bob",
                     "Description for flag b|bob flag");
    root.addShortFlagWithParameter("c",
                                   "Description for c flag",
                                   "word",
                                   "Some word, whatever you want",
                                   CliArgumentType.STRING);
    root.addLongFlagWithParameter("d",
                                  "dog",
                                  "Description for d|dog",
                                  "name",
                                  "The name of the dog",
                                  CliArgumentType.STRING);
    root.addShortFlagWithParameter("e",
                                   "This is the e short flag",
                                   "number",
                                   "This is your favorite number",
                                   CliArgumentType.INTEGER);

    CliList tunaList = root.addList("tuna",
                                    "This is the tuna command list");
    tunaList.addLongFlagWithParameter("a",
                                      "andrew",
                                      "Description for andrew flag",
                                      "age",
                                      "The age you think I am",
                                      CliArgumentType.INTEGER);
    tunaList.addShortFlag("f",
                          "The f flag, ya know");

    CliCommand tunaMarlinCommand = tunaList.addCommand("marlin",
                                                      "This is the marlin command",
                                                      tunaMarlinAction);
    tunaMarlinCommand.addShortFlag("z",
                                   "The z flag, ya know");

    root.addCommand("mayo", "This is the mayo command", mayoAction);
    return cli;
  }

  /*
   * Section - Flag Tests
   */

  /*
   * Subsection - Negative Flag Tests
   */

  @Test(expected = IllegalArgumentException.class)
  public void testBadFlagSyntaxStart() throws IllegalArgumentException {
    parse("---a", "-b");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadFlagSyntaxEnd() throws IllegalArgumentException {
    parse("-b", "---a");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testUnknownShortFlag() throws IllegalArgumentException {
    parse("-b", "-z");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testUnknownLongFlag() throws IllegalArgumentException {
    parse("-b", "--zebra");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testNoShortArgument() throws IllegalArgumentException {
    parse("-b", "-c");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testNoLongArgument() throws IllegalArgumentException {
    parse("-b", "--dog");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testLongFlagShortSyntax() throws IllegalArgumentException {
    parse("-b", "-bob");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadCommand() throws IllegalArgumentException {
    parse("-a", "this-command-does-not-exist");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadNestedCommand() throws IllegalArgumentException {
    parse("-a", "tuna", "this-command-does-not-exist");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadFlagTypeAtEnd() throws IllegalArgumentException {
    parse("-a", "--dog", "rover", "-e", "this is not a number");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadFlagTypeAtBeginning() throws IllegalArgumentException {
    parse("-e", "moooo", "-a", "--dog", "rover");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadGetFlagValue() throws IllegalArgumentException {
    parse("-e", "15", "tuna", "marlin");
    assertTrue(tunaMarlinAction.getRan());
    tunaMarlinAction.getFlags().getValue("e", CliArgumentType.STRING);
  }

  /*
   * Subsection - Positive Flag Tests
   */

  @Test
  public void shortFlagOnlyTest() throws IllegalArgumentException {
    parse("-a");
  }

  @Test
  public void longFlagShortFlagTest() throws IllegalArgumentException {
    parse("--bob", "-a");
  }

  @Test
  public void testShortArgumentShortFlag() throws IllegalArgumentException {
    parse("-c", "hello", "-a");
  }

  @Test
  public void testEverything() throws IllegalArgumentException {
    parse("-c", "hello", "--bob", "-a", "--dog", "world");
  }

  @Test
  public void testEmptyArgs() {
    parse();
  }

  /*
   * Section - Commands
   */

  /*
   * Subsection - Negative Commands
   */

  @Test(expected = IllegalArgumentException.class)
  public void testTunaListWithArguments() {
    parse("tuna", "hello", "world");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testTunaListWithParameterBeforeCommand() {
    parse("tuna", "hello", "marlin");
  }

  /*
   * Subsection - Positive Commands
   */

  @Test
  public void testTunaListWithoutPreFlags() {
    parse("tuna");
    assertFalse(tunaMarlinAction.getRan());
  }

  @Test
  public void testTunaListWithPreFlags() {
    parse("-a", "-c", "hello", "tuna");
    assertFalse(tunaMarlinAction.getRan());
  }

  @Test
  public void testMarlinCommandWithoutArgument() {
    parse("tuna", "marlin");
    assertTrue(tunaMarlinAction.getRan());
    CliFlags flags = tunaMarlinAction.getFlags();
    assertArrayEquals(new String[0], flags.getAllShortFlags());
    assertNull(flags.getValue("z", CliArgumentType.BOOLEAN));
    assertEquals(0, tunaMarlinAction.getArguments().length);
  }

  @Test
  public void testMarlinCommandWithFlag() {
    parse("tuna", "marlin", "-z");
    assertTrue(tunaMarlinAction.getRan());
    CliFlags flags = tunaMarlinAction.getFlags();
    assertArrayEquals(new String[] { "z" }, flags.getAllShortFlags());
    assertEquals(Boolean.TRUE, flags.getValue("z", CliArgumentType.BOOLEAN));
    assertArrayEquals(tunaMarlinAction.getArguments(), new String[0]);
  }

  @Test
  public void testMarlinCommandWithArguments() {
    parse("tuna", "marlin", "hello", "world");
    assertTrue(tunaMarlinAction.getRan());
    CliFlags flags = tunaMarlinAction.getFlags();
    assertArrayEquals(flags.getAllShortFlags(), new String[0]);
    assertNull(flags.getValue("z", CliArgumentType.BOOLEAN));
    assertArrayEquals(new String[] { "hello", "world" }, tunaMarlinAction.getArguments());
  }

  @Test
  public void testMarlinCommandWithArgumentsAndFlag() {
    parse("tuna", "marlin", "-z", "hello", "world");
    CliFlags flags = tunaMarlinAction.getFlags();
    assertArrayEquals(flags.getAllShortFlags(), new String[] { "z" });
    assertEquals(Boolean.TRUE, flags.getValue("z", CliArgumentType.BOOLEAN));
    assertArrayEquals(new String[] { "hello", "world" }, tunaMarlinAction.getArguments());
  }

  @Test
  public void testMarlinCommandWitPreFlagsAndArguments() {
    parse("tuna", "-a", "15", "-f", "marlin", "hello", "world");
    assertTrue(tunaMarlinAction.getRan());
    CliFlags flags = tunaMarlinAction.getFlags();
    assertArrayEquals(new String[] { "a", "f" }, flags.getAllShortFlags());
    assertEquals(15, flags.getValue("a", CliArgumentType.INTEGER));
  }

  @Test
  public void testMayoCommand() {
    parse("mayo", "a", "b", "c");
    assertTrue(mayoAction.getRan());
    CliFlags flags = mayoAction.getFlags();
    assertArrayEquals(new String[0], flags.getAllShortFlags());
    assertArrayEquals(new String[] { "a", "b", "c" }, mayoAction.getArguments());
  }

  /*
   * Section - Visitor
   */
  @Test
  public void testVisitorPattern() {
    TestCliVisitor visitor = new TestCliVisitor();
    visit(visitor);

    assertVisited(visitor.getVisitedShortFlags(), "a", "f", "z");
    assertVisited(visitor.getVisitedShortFlagsWithParameters(), "c", "e");
    assertVisited(visitor.getVisitedLongFlags(), "b");
    assertVisited(visitor.getVisitedLongFlagsWithParameters(), "d", "a");
    assertVisited(visitor.getVisitedCommands(), "mayo", "marlin");
    assertVisited(visitor.getVisitedLists(), "cli-test", "tuna");
    assertVisited(visitor.getLeftLists(), "cli-test", "tuna");
  }

  private void assertVisited(List<String> flags, String...expecteds) {
    assertEquals(flags.size(), expecteds.length);
    for (String expected : expecteds) {
      assertTrue(flags.contains(expected));
    }
  }
}
