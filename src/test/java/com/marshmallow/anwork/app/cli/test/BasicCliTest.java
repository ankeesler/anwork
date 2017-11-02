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
import com.marshmallow.anwork.core.test.TestUtilities;

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
  private TestCliAction fooAction = new TestCliAction();

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
        .setArgument("word", ArgumentType.STRING).setDescription("Some word, whatever you want");
    root.addFlag("d")
        .setLongFlag("dog")
        .setDescription("Description for d|dog")
        .setArgument("name", ArgumentType.STRING).setDescription("The name of the dog");
    root.addFlag("e")
        .setDescription("This is the e short flag")
        .setArgument("number", ArgumentType.NUMBER).setDescription("This is your favorite number");

    // Change the name of the list here to make sure that the changed list name will be parsed.
    MutableList tunaList = root.addList("wrong")
                               .setDescription("This is the tuna command list")
                               .setName("tuna");
    tunaList.addFlag("w")
            .setLongFlag("andrew")
            .setArgument("age", ArgumentType.NUMBER).setDescription("The age you think I am");
    tunaList.addFlag("f")
            .setDescription("The f flag, ya know");

    MutableCommand tunaMarlinCommand = tunaList.addCommand("marlin", tunaMarlinAction);
    tunaMarlinCommand.setDescription("This is the marlin command");
    tunaMarlinCommand.addFlag("wrong").setDescription("The z flag, ya know").setShortFlag("z");
    // ^^^ this makes sure that we can change the short flag on a flag and parse the correct flag
    tunaMarlinCommand.addFlag("t").setLongFlag("donut");
    tunaMarlinCommand.addFlag("t").setLongFlag("duplicate");
    tunaMarlinCommand.addArgument("number", ArgumentType.NUMBER);
    tunaMarlinCommand.addArgument("wrong", ArgumentType.STRING)
                     .setName("name")
                     .setDescription("hey");
    // ^^^ this makes sure that we can change the name of an argument when building the CLI

    // Change the name of the command here to make sure the changed command name will be parsed.
    root.addCommand("wrong", mayoAction)
        .setDescription("This is the mayo command")
        .setName("mayo");

    root.addCommand("command-without-description", new TestCliAction());
    MutableList duplicateList = root.addList("duplicate-list")
                                    .setDescription("this has a description!");
    duplicateList = root.addList("duplicate-list"); // no description...
    duplicateList.addCommand("foo", new TestCliAction());
    duplicateList.addCommand("foo", fooAction);

    root.addList("list-without-description");
    root.addFlag("f");
    // ^^^ this flag ensures that we can add the same flag to different paths in the CLI tree

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
    parse("-e", "15", "tuna", "marlin", "25", "steve");
    assertTrue(tunaMarlinAction.getRan());
    assertTrue(tunaMarlinAction.getFlags().containsKey("e"));
    tunaMarlinAction.getFlags().getValue("e", ArgumentType.STRING);
  }

  @Test(expected = IllegalArgumentException.class)
  public void testDuplicateFlagsNegative() {
    parse("tuna", "marlin", "--donut", "25", "steve");
  }

  @Test(expected = IllegalStateException.class)
  public void testAssigningSameFlagToCliTreePath() {
    Cli cli = new Cli("foo");
    MutableList rootList = cli.getRoot();
    rootList.addFlag("a");
    rootList.addCommand("bar", (flags, args) -> { }).addFlag("a");
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

  @Test
  public void testDuplicateFlagsPositive() {
    parse("tuna", "marlin", "--duplicate", "25", "steve");
    assertTrue(tunaMarlinAction.getRan());
    assertTrue(tunaMarlinAction.getFlags().containsKey("t"));
    Boolean b = tunaMarlinAction.getFlags().getValue("t", ArgumentType.BOOLEAN);
    assertEquals(Boolean.TRUE, b);
  }

  @Test
  public void testNonexistentFlagValues() {
    parse("tuna", "marlin", "-t", "25", "steve");
    assertFalse(tunaMarlinAction.getFlags().containsKey(""));
    assertFalse(tunaMarlinAction.getFlags().containsKey("a"));
    assertFalse(tunaMarlinAction.getFlags().containsKey("b"));
    assertFalse(tunaMarlinAction.getFlags().containsKey("c"));
    assertFalse(tunaMarlinAction.getFlags().containsKey("e"));
    assertFalse(tunaMarlinAction.getFlags().containsKey("f"));
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

  @Test(expected = IllegalArgumentException.class)
  public void testMarlinCommandMissingArguments() {
    parse("tuna", "marlin");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testMarlinCommandWithFlagButMissingArguments() {
    parse("tuna", "marlin", "-z");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testMarlinCommandWithTooFewArguments() {
    parse("tuna", "marlin", "25");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testMarlinCommandWithTooManyArguments() {
    parse("tuna", "marlin", "25", "steve", "bacon");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testMarlinCommandWithBadArgumentTypes() {
    parse("tuna", "marlin", "steve", "25");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testMayoCommandWithArguments() {
    parse("mayo", "hey");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testGettingWrongArgumentType() {
    parse("tuna", "marlin", "25", "steve");
    ArgumentValues arguments = tunaMarlinAction.getArguments();
    arguments.getValue("number", ArgumentType.STRING);
  }

  @Test(expected = IllegalArgumentException.class)
  public void testUnkownRootedCommand() {
    parse("moo");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testUnkownNestedCommand() {
    parse("tuna", "marlin");
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
  public void testMarlinCommandWithArguments() {
    parse("tuna", "marlin", "25", "steve");
    assertTrue(tunaMarlinAction.getRan());

    ArgumentValues flags = tunaMarlinAction.getFlags();
    assertArrayEquals(flags.getAllKeys(), new String[0]);
    assertNull(flags.getValue("z", ArgumentType.BOOLEAN));

    ArgumentValues arguments = tunaMarlinAction.getArguments();
    assertEquals(2, arguments.getAllKeys().length);
    assertTrue(arguments.containsKey("number"));
    assertTrue(arguments.containsKey("name"));
    assertEquals(new Long(25), arguments.getValue("number", ArgumentType.NUMBER));
    assertEquals("steve", arguments.getValue("name", ArgumentType.STRING));
    assertNull(arguments.getValue("wrong", ArgumentType.STRING));
  }

  @Test
  public void testMarlinCommandWithArgumentsAndFlag() {
    parse("tuna", "marlin", "-z", "25", "steve");
    assertTrue(tunaMarlinAction.getRan());

    ArgumentValues flags = tunaMarlinAction.getFlags();
    assertArrayEquals(flags.getAllKeys(), new String[] { "z" });
    assertEquals(Boolean.TRUE, flags.getValue("z", ArgumentType.BOOLEAN));

    ArgumentValues arguments = tunaMarlinAction.getArguments();
    assertEquals(2, arguments.getAllKeys().length);
    assertTrue(arguments.containsKey("number"));
    assertTrue(arguments.containsKey("name"));
    assertEquals(new Long(25), arguments.getValue("number", ArgumentType.NUMBER));
    assertEquals("steve", arguments.getValue("name", ArgumentType.STRING));
    assertNull(arguments.getValue("wrong", ArgumentType.STRING));
  }

  @Test
  public void testMarlinCommandWitPreFlagsAndArguments() {
    parse("tuna", "-w", "15", "-f", "marlin", "25", "steve");
    assertTrue(tunaMarlinAction.getRan());

    ArgumentValues flags = tunaMarlinAction.getFlags();
    assertArrayEquals(new String[] { "f", "w" }, flags.getAllKeys());
    Long number = flags.getValue("w", ArgumentType.NUMBER);
    assertEquals(new Long(15), number);

    ArgumentValues arguments = tunaMarlinAction.getArguments();
    assertEquals(2, arguments.getAllKeys().length);
    assertTrue(arguments.containsKey("number"));
    assertTrue(arguments.containsKey("name"));
    assertEquals(new Long(25), arguments.getValue("number", ArgumentType.NUMBER));
    assertEquals("steve", arguments.getValue("name", ArgumentType.STRING));
    assertNull(arguments.getValue("wrong", ArgumentType.STRING));
  }

  @Test
  public void testMayoCommand() {
    parse("mayo");
    assertTrue(mayoAction.getRan());
    ArgumentValues flags = mayoAction.getFlags();
    assertArrayEquals(new String[0], flags.getAllKeys());
    ArgumentValues arguments = mayoAction.getArguments();
    assertEquals(0, arguments.getAllKeys().length);
  }

  @Test
  public void testDuplicateList() {
    parse("duplicate-list");
  }

  @Test
  public void testDuplicateCommand() {
    parse("duplicate-list", "foo");
    assertTrue(fooAction.getRan());
  }

  /*
   * Section - Visitor
   */

  @Test
  public void testVisitorPattern() {
    TestCliVisitor visitor = new TestCliVisitor();
    visit(visitor);
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedShortFlags(),
                                            "a", "f", "f", "z");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedShortFlagsWithParameters(),
                                            "c", "e");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedLongFlags(),
                                            "bob", "duplicate");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedLongFlagsWithParameters(),
                                            "dog", "andrew");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedCommands(),
                                            "command-without-description", "mayo",
                                            "foo", "marlin");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedCommandArguments(),
                                            "number", "name");
    TestUtilities.assertVariadicArrayEquals(visitor.getVisitedLists(),
                                            "cli-test", "duplicate-list",
                                            "list-without-description", "tuna");
    TestUtilities.assertVariadicArrayEquals(visitor.getLeftLists(),
                                            "duplicate-list", "list-without-description",
                                            "tuna", "cli-test");
  }

  @Test
  public void testHasXxxMethods() {
    // This test makes sure that if we don't set an optional field, then the hasXxx methods will
    // return false.
    OptionalDataCliVisitor visitor = new OptionalDataCliVisitor();
    visit(visitor);
    TestUtilities.assertVariadicArrayEquals(visitor.getFlagsWithDescriptions(),
                                            "a", "b", "d", "e", "f", "z");
    TestUtilities.assertVariadicArrayEquals(visitor.getCommandsWithDescriptions(),
                                            "mayo", "marlin");
    TestUtilities.assertVariadicArrayEquals(visitor.getCommandArgumentsWithDescriptions(),
                                            "name");
    TestUtilities.assertVariadicArrayEquals(visitor.getListsWithDescriptions(),
                                            "cli-test", "tuna");
  }
}
