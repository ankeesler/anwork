package com.marshmallow.anwork.app.cli.test;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliCommand;
import com.marshmallow.anwork.app.cli.CliList;

import org.junit.Before;
import org.junit.Test;

/**
 * This is a test for the CLI.
 *
 * @author Andrew
 * Created Sep 9, 2017
 */
public class CliTest {

  private Cli cli = new Cli("cli-test",  "The are commands for the CLI unit test");
  private TestCliAction aShortFlagAction = new TestCliAction();
  private TestCliAction bLongFlagAction = new TestCliAction();
  private TestCliAction cShortParameterAction = new TestCliAction();
  private TestCliAction dLongParameterAction = new TestCliAction();

  private TestCliAction tunaAndrewParameterAction = new TestCliAction();
  private TestCliAction tunaFShortFlagAction = new TestCliAction();

  private TestCliAction tunaMarlinAction = new TestCliAction();
  private TestCliAction tunaMarlinZShortFlagAction = new TestCliAction();

  private TestCliAction mayoAction = new TestCliAction();

  @Before
  public void setupCli() {
    CliList root = cli.getRoot();
    root.addShortFlag("a",
                      "Description for a flag",
                      aShortFlagAction);
    root.addLongFlag("b",
                     "bob",
                     "Description for flag b|bob flag",
                     bLongFlagAction);
    root.addShortFlagWithParameter("c",
                                   "Description for c flag",
                                   "word",
                                   cShortParameterAction);
    root.addLongFlagWithParameter("d",
                                  "dog",
                                  "Description for d|dog",
                                  "name",
                                  dLongParameterAction);

    CliList tunaList = root.addList("tuna",
                                    "This is the tuna command list");
    tunaList.addLongFlagWithParameter("a",
                                      "andrew",
                                      "Description for andrew flag",
                                      "whatever",
                                      tunaAndrewParameterAction);
    tunaList.addShortFlag("f",
                          "The f flag, ya know",
                          tunaFShortFlagAction);

    CliCommand tunaMarlinCommand = tunaList.addCommand("marlin",
                                                      "This is the marlin command",
                                                      tunaMarlinAction);
    tunaMarlinCommand.addShortFlag("z",
                                   "The z flag, ya know",
                                   tunaMarlinZShortFlagAction);

    root.addCommand("mayo", "This is the mayo command", mayoAction);
  }

  private void runTest(String...args) {
    cli.parse(args);
  }

  /*
   * Section - Flag Tests
   */

  /*
   * Subsection - Negative Flag Tests
   */

  @Test(expected = IllegalArgumentException.class)
  public void testBadFlagSyntaxStart() throws IllegalArgumentException {
    runTest("---a", "-b" );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadFlagSyntaxEnd() throws IllegalArgumentException {
    runTest("-b", "---a" );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testUnknownShortFlag() throws IllegalArgumentException {
    runTest("-b", "-z");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testUnknownLongFlag() throws IllegalArgumentException {
    runTest("-b", "--zebra");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testNoShortArgument() throws IllegalArgumentException {
    runTest("-b", "-c");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testNoLongArgument() throws IllegalArgumentException {
    runTest("-b", "--dog");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testLongFlagShortSyntax() throws IllegalArgumentException {
    runTest("-b", "-bob");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadCommand() throws IllegalArgumentException {
    runTest("-a", "this-command-does-not-exist");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadNestedCommand() throws IllegalArgumentException {
    runTest("-a", "tuna", "this-command-does-not-exist");
  }

  /*
   * Subsection - Positive Flag Tests
   */

  @Test
  public void shortFlagOnlyTest() throws IllegalArgumentException {
    runTest("-a");
    CliTestUtilities.assertActionRan(aShortFlagAction);
    CliTestUtilities.assertActionDidNotRun(bLongFlagAction);
    CliTestUtilities.assertActionDidNotRun(cShortParameterAction);
    CliTestUtilities.assertActionDidNotRun(dLongParameterAction);
  }

  @Test
  public void longFlagShortFlagTest() throws IllegalArgumentException {
    runTest("--bob", "-a");
    CliTestUtilities.assertActionRan(aShortFlagAction);
    CliTestUtilities.assertActionRan(bLongFlagAction);
    CliTestUtilities.assertActionDidNotRun(cShortParameterAction);
    CliTestUtilities.assertActionDidNotRun(dLongParameterAction);
  }

  @Test
  public void testShortArgumentShortFlag() throws IllegalArgumentException {
    runTest("-c", "hello", "-a");
    CliTestUtilities.assertActionRan(aShortFlagAction);
    CliTestUtilities.assertActionDidNotRun(bLongFlagAction);
    CliTestUtilities.assertActionRan(cShortParameterAction, "hello");
    CliTestUtilities.assertActionDidNotRun(dLongParameterAction);
  }

  @Test
  public void testEverything() throws IllegalArgumentException {
    runTest("-c", "hello", "--bob", "-a", "--dog", "world");
    CliTestUtilities.assertActionRan(aShortFlagAction);
    CliTestUtilities.assertActionRan(bLongFlagAction);
    CliTestUtilities.assertActionRan(cShortParameterAction, "hello");
    CliTestUtilities.assertActionRan(dLongParameterAction, "world");
  }

  @Test
  public void testEmptyArgs() {
    runTest();
    CliTestUtilities.assertActionDidNotRun(aShortFlagAction);
    CliTestUtilities.assertActionDidNotRun(bLongFlagAction);
    CliTestUtilities.assertActionDidNotRun(cShortParameterAction);
    CliTestUtilities.assertActionDidNotRun(dLongParameterAction);
  }

  /*
   * Section - Commands
   */

  /*
   * Subsection - Negative Commands
   */

  @Test(expected = IllegalArgumentException.class)
  public void testTunaListWithArguments() {
    runTest("tuna", "hello", "world");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testTunaListWithParameterBeforeCommand() {
    runTest("tuna", "hello", "marlin");
  }

  /*
   * Subsection - Positive Commands
   */

  @Test
  public void testTunaListWithoutPreFlags() {
    runTest("tuna");
    CliTestUtilities.assertActionDidNotRun(aShortFlagAction);
    CliTestUtilities.assertActionDidNotRun(cShortParameterAction);
  }

  @Test
  public void testTunaListWithPreFlags() {
    runTest("-a", "-c", "hello", "tuna");
    CliTestUtilities.assertActionRan(aShortFlagAction);
    CliTestUtilities.assertActionRan(cShortParameterAction, "hello");
  }

  @Test
  public void testMarlinCommandWithoutArgument() {
    runTest("tuna", "marlin");
    CliTestUtilities.assertActionRan(tunaMarlinAction);
  }

  @Test
  public void testMarlinCommandWithFlag() {
    runTest("tuna", "marlin", "-z");
    CliTestUtilities.assertActionRan(tunaMarlinAction);
    CliTestUtilities.assertActionRan(tunaMarlinZShortFlagAction);
  }

  @Test
  public void testMarlinCommandWithArguments() {
    runTest("tuna", "marlin", "hello", "world");
    CliTestUtilities.assertActionRan(tunaMarlinAction, "hello" , "world");
    CliTestUtilities.assertActionDidNotRun(tunaMarlinZShortFlagAction);
  }

  @Test
  public void testMarlinCommandWithArgumentsAndFlag() {
    runTest("tuna", "marlin", "-z", "hello", "world");
    CliTestUtilities.assertActionRan(tunaMarlinAction, "hello", "world");
    CliTestUtilities.assertActionRan(tunaMarlinZShortFlagAction);
  }

  @Test
  public void testMarlinCommandWitPreFlagsAndArguments() {
    runTest("tuna", "-a", "andrew", "-f", "marlin", "hello", "world");
    CliTestUtilities.assertActionRan(tunaAndrewParameterAction, "andrew");
    CliTestUtilities.assertActionRan(tunaFShortFlagAction);
    CliTestUtilities.assertActionRan(tunaMarlinAction, "hello", "world");
  }

  @Test
  public void usageTest() {
    System.out.println(cli.getUsage());
  }
}
