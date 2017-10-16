package com.marshmallow.anwork.app.cli.test;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.assertTrue;

import com.marshmallow.anwork.app.cli.ArgumentType;
import com.marshmallow.anwork.app.cli.ArgumentValues;
import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.MutableCommand;
import com.marshmallow.anwork.app.cli.MutableList;

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

  // TODO: the description on the root node is correct
  // TODO: changing the name of a node doesn't mess anything up
  // TODO: changing the name of a short flag doesn't mess anything up
  // TODO: not setting a description doesn't throw any errors
  // TODO: test the #hasXXX functionality
  // TODO: test visitor functionality more robustly
  // TODO: get rid of warnings in this file
  // TODO: update Cli javadoc
  // TODO: update CLI XML
  // TODO: update CLI ARCH.md entry

  private TestCliAction tunaMarlinAction = new TestCliAction();

  private TestCliAction mayoAction = new TestCliAction();

  @Override
  protected Cli createCli() {
    Cli cli = new Cli("cli-test");
    MutableList root = cli.getRoot();
    root.setDescription("The are commands for the CLI unit test");
    root.addFlag("a")
        .setDescription("Description for a flag");
    root.addFlag("b")
        .setLongFlag("bob")
        .setDescription("Description for flag b|bob flag");
    root.addFlag("c")
        .setDescription("Description for c flag")
        .setArgument("word", ArgumentType.STRING).setDescription("Some word, whatever you want");
    root.addFlag("d")
        .setLongFlag("dog")
        .setDescription("Description for d|dog")
        .setArgument("name", ArgumentType.STRING).setDescription("The name of the dog");
    root.addFlag("e")
        .setDescription("This is the e short flag")
        .setArgument("number", ArgumentType.NUMBER).setDescription("This is your favorite number");

    MutableList tunaList = root.addList("tuna").setDescription("This is the tuna command list");
    tunaList.addFlag("a")
            .setLongFlag("andrew")
            .setDescription("Description for andrew flag")
            .setArgument("age", ArgumentType.NUMBER).setDescription("The age you think I am");
    tunaList.addFlag("f")
            .setDescription("The f flag, ya know");

    MutableCommand tunaMarlinCommand = tunaList.addCommand("marlin", tunaMarlinAction);
    tunaMarlinCommand.setDescription("This is the marlin command");
    tunaMarlinCommand.addFlag("z").setDescription("The z flag, ya know");

    root.addCommand("mayo", mayoAction).setDescription("This is the mayo command");

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
    tunaMarlinAction.getFlags().getValue("e", ArgumentType.STRING);
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
    ArgumentValues flags = tunaMarlinAction.getFlags();
    assertArrayEquals(new String[0], flags.getAllKeys());
    assertNull(flags.getValue("z", ArgumentType.BOOLEAN));
    assertEquals(0, tunaMarlinAction.getArguments().length);
  }

  @Test
  public void testMarlinCommandWithFlag() {
    parse("tuna", "marlin", "-z");
    assertTrue(tunaMarlinAction.getRan());
    ArgumentValues flags = tunaMarlinAction.getFlags();
    assertArrayEquals(new String[] { "z" }, flags.getAllKeys());
    assertEquals(Boolean.TRUE, flags.getValue("z", ArgumentType.BOOLEAN));
    assertArrayEquals(tunaMarlinAction.getArguments(), new String[0]);
  }

  @Test
  public void testMarlinCommandWithArguments() {
    parse("tuna", "marlin", "hello", "world");
    assertTrue(tunaMarlinAction.getRan());
    ArgumentValues flags = tunaMarlinAction.getFlags();
    assertArrayEquals(flags.getAllKeys(), new String[0]);
    assertNull(flags.getValue("z", ArgumentType.BOOLEAN));
    assertArrayEquals(new String[] { "hello", "world" }, tunaMarlinAction.getArguments());
  }

  @Test
  public void testMarlinCommandWithArgumentsAndFlag() {
    parse("tuna", "marlin", "-z", "hello", "world");
    ArgumentValues flags = tunaMarlinAction.getFlags();
    assertArrayEquals(flags.getAllKeys(), new String[] { "z" });
    assertEquals(Boolean.TRUE, flags.getValue("z", ArgumentType.BOOLEAN));
    assertArrayEquals(new String[] { "hello", "world" }, tunaMarlinAction.getArguments());
  }

  @Test
  public void testMarlinCommandWitPreFlagsAndArguments() {
    parse("tuna", "-a", "15", "-f", "marlin", "hello", "world");
    assertTrue(tunaMarlinAction.getRan());
    ArgumentValues flags = tunaMarlinAction.getFlags();
    assertArrayEquals(new String[] { "a", "f" }, flags.getAllKeys());
    Long number = flags.getValue("a", ArgumentType.NUMBER);
    assertEquals(new Long(15), number);
  }

  @Test
  public void testMayoCommand() {
    parse("mayo", "a", "b", "c");
    assertTrue(mayoAction.getRan());
    ArgumentValues flags = mayoAction.getFlags();
    assertArrayEquals(new String[0], flags.getAllKeys());
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
    assertVisited(visitor.getVisitedLongFlags(), "bob");
    assertVisited(visitor.getVisitedLongFlagsWithParameters(), "dog", "andrew");
    assertVisited(visitor.getVisitedCommands(), "mayo", "marlin");
    assertVisited(visitor.getVisitedLists(), "cli-test", "tuna");
    assertVisited(visitor.getLeftLists(), "tuna", "cli-test");
  }

  private void assertVisited(String[] visitedStuff, String...expecteds) {
    assertArrayEquals(expecteds, visitedStuff);
  }
}
