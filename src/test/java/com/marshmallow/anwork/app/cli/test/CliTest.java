package com.marshmallow.anwork.app.cli.test;

import org.junit.Test;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliFlag;
import com.marshmallow.anwork.app.cli.CliNode;

import static org.junit.Assert.*;

import org.junit.Before;

/**
 * This is a test for the CLI.
 *
 * @author Andrew
 * @date Sep 9, 2017
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

  private TestCliAction mayoAction = new TestCliAction();

  @Before
  public void setupCli() {
    cli.addFlag(CliFlag.makeShortFlag("a", "Description for a flag", aShortFlagAction));
    cli.addFlag(CliFlag.makeLongFlag("b", "bob", "Description for flag b|bob flag", bLongFlagAction));
    cli.addFlag(CliFlag.makeShortFlagWithParameter("c", "Description for c flag", "word", cShortParameterAction));
    cli.addFlag(CliFlag.makeLongFlagWithParameter("d", "dog", "Description for d|dog", "name", dLongParameterAction));

    CliNode tunaList = CliNode.makeList("tuna", "This is the tuna command");
    tunaList.addFlag(CliFlag.makeLongFlagWithParameter("a", "andrew", "Description for andrew flag", "whatever", tunaAndrewParameterAction));
    tunaList.addFlag(CliFlag.makeShortFlag("f", "The f flag, ya know", tunaFShortFlagAction));
    cli.addNode(tunaList);

    CliNode tunaMarlinCommand = CliNode.makeCommand("marlin", "This is the marlin command", tunaMarlinAction);
    tunaList.addNode(tunaMarlinCommand);

    CliNode mayoCommand = CliNode.makeCommand("mayo", "This is the mayo command", mayoAction); 
    cli.addNode(mayoCommand);
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadFlagSyntaxStart() throws IllegalArgumentException {
    cli.parse(new String[] { "---a", "-b" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadFlagSyntaxEnd() throws IllegalArgumentException {
    cli.parse(new String[] { "-b", "---a" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testUnknownShortFlag() throws IllegalArgumentException {
    cli.parse(new String[] { "-b", "-z" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testUnknownLongFlag() throws IllegalArgumentException {
    cli.parse(new String[] { "-b", "--zebra" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testNoShortArgument() throws IllegalArgumentException {
    cli.parse(new String[] { "-b", "-c" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testNoLongArgument() throws IllegalArgumentException {
    cli.parse(new String[] { "-b", "--dog" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testLongFlagShortSyntax() throws IllegalArgumentException {
    cli.parse(new String[] { "-b", "-bob" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadCommand() throws IllegalArgumentException {
    cli.parse(new String[] { "-a", "this-command-does-not-exist" });
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadNestedCommand() throws IllegalArgumentException {
    cli.parse(new String[] { "-a", "tuna", "this-command-does-not-exist" });
  }

  @Test
  public void shortFlagOnlyTest() throws IllegalArgumentException {
    cli.parse(new String[] { "-a" });
    assertTrue(aShortFlagAction.getRan());
    CliTestUtilities.assertActionRan(aShortFlagAction);
    assertFalse(bLongFlagAction.getRan());
    assertFalse(cShortParameterAction.getRan());
    assertFalse(dLongParameterAction.getRan());
  }

  @Test
  public void longFlagShortFlagTest() throws IllegalArgumentException {
    cli.parse(new String[] { "--bob", "-a" });
    assertTrue(aShortFlagAction.getRan());
    CliTestUtilities.assertActionRan(aShortFlagAction);
    assertTrue(bLongFlagAction.getRan());
    CliTestUtilities.assertActionRan(bLongFlagAction);
    assertFalse(cShortParameterAction.getRan());
    assertFalse(dLongParameterAction.getRan());
  }

  @Test
  public void testShortArgumentShortFlag() throws IllegalArgumentException {
    cli.parse(new String[] { "-c", "hello", "-a" });
    assertTrue(aShortFlagAction.getRan());
    CliTestUtilities.assertActionRan(aShortFlagAction);
    assertFalse(bLongFlagAction.getRan());
    assertTrue(cShortParameterAction.getRan());
    CliTestUtilities.assertActionRan(cShortParameterAction, "hello");
    assertFalse(dLongParameterAction.getRan());
  }

  @Test
  public void testEverything() throws IllegalArgumentException {
    cli.parse(new String[] { "-c", "hello", "--bob", "-a", "--dog", "world" });
    assertTrue(aShortFlagAction.getRan());
    CliTestUtilities.assertActionRan(aShortFlagAction);
    assertTrue(bLongFlagAction.getRan());
    CliTestUtilities.assertActionRan(bLongFlagAction);
    assertTrue(cShortParameterAction.getRan());
    CliTestUtilities.assertActionRan(cShortParameterAction, "hello");
    assertTrue(dLongParameterAction.getRan());
    CliTestUtilities.assertActionRan(dLongParameterAction, "world");
  }

  @Test
  public void testEmptyArgs() {
    cli.parse(new String[] { });
    assertFalse(aShortFlagAction.getRan());
    assertFalse(bLongFlagAction.getRan());
    assertFalse(cShortParameterAction.getRan());
    assertFalse(dLongParameterAction.getRan());
  }

  @Test
  public void testTunaListWithoutPreFlags() {
    cli.parse(new String[] { "tuna" } );
    assertFalse(aShortFlagAction.getRan());
    assertFalse(cShortParameterAction.getRan());
  }

  @Test
  public void testTunaListWithPreFlags() {
    cli.parse(new String[] { "-a", "-c", "hello", "tuna" } );
    CliTestUtilities.assertActionRan(aShortFlagAction);
    CliTestUtilities.assertActionRan(cShortParameterAction, "hello");
  }

  @Test(expected = IllegalArgumentException.class)
  public void testTunaListWithArguments() {
    cli.parse(new String[] { "tuna", "hello", "world" } );
  }

  @Test
  public void testMarlinCommandWithoutArgument() {
    cli.parse(new String[] { "tuna", "marlin" } );
    CliTestUtilities.assertActionRan(tunaMarlinAction);
  }

  @Test
  public void testMarlinCommandWithArguments() {
    cli.parse(new String[] { "tuna", "marlin", "hello", "world" } );
    CliTestUtilities.assertActionRan(tunaMarlinAction, "hello" , "world");
  }

  @Test
  public void testMarlinCommandWitPreFlagsAndArguments() {
    cli.parse(new String[] { "tuna", "-a", "andrew", "-f", "marlin", "hello", "world" } );
    CliTestUtilities.assertActionRan(tunaAndrewParameterAction, "andrew");
    CliTestUtilities.assertActionRan(tunaFShortFlagAction);
    CliTestUtilities.assertActionRan(tunaMarlinAction, "hello" , "world");
  }

  @Test
  public void usageTest() {
    System.out.println(cli.getUsage());
  }
}
