package com.marshmallow.anwork.app.cli.test;

import org.junit.Test;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliFlag;

import static org.junit.Assert.*;

import org.junit.Before;

/**
 * This is a test for the CLI.
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public class CliTest {

  private Cli cli = new Cli();
  private TestCliAction aShortFlagAction = new TestCliAction();
  private TestCliAction bLongFlagAction = new TestCliAction();
  private TestCliAction cShortParameterAction = new TestCliAction();
  private TestCliAction dLongParameterAction = new TestCliAction();

  @Before
  public void setupCli() {
    cli.addFlag(CliFlag.makeShortFlag("a", "Description for a flag", aShortFlagAction));
    cli.addFlag(CliFlag.makeLongFlag("b", "bob", "Description for flag b|bob flag", bLongFlagAction));
    cli.addFlag(CliFlag.makeShortFlagWithParameter("c", "Description for c flag", "word", cShortParameterAction));
    cli.addFlag(CliFlag.makeLongFlagWithParameter("d", "dog", "Description for d|dog", "name", dLongParameterAction));
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

  @Test
  public void shortFlagOnlyTest() throws IllegalArgumentException {
    cli.parse(new String[] { "-a" });
    assertTrue(aShortFlagAction.getRan());
    CliTestUtilities.assertValidArguments(aShortFlagAction);
    assertFalse(bLongFlagAction.getRan());
    assertFalse(cShortParameterAction.getRan());
    assertFalse(dLongParameterAction.getRan());
  }

  @Test
  public void longFlagShortFlagTest() throws IllegalArgumentException {
    cli.parse(new String[] { "--bob", "-a" });
    assertTrue(aShortFlagAction.getRan());
    CliTestUtilities.assertValidArguments(aShortFlagAction);
    assertTrue(bLongFlagAction.getRan());
    CliTestUtilities.assertValidArguments(bLongFlagAction);
    assertFalse(cShortParameterAction.getRan());
    assertFalse(dLongParameterAction.getRan());
  }

  @Test
  public void testShortArgumentShortFlag() throws IllegalArgumentException {
    cli.parse(new String[] { "-c", "hello", "-a" });
    assertTrue(aShortFlagAction.getRan());
    CliTestUtilities.assertValidArguments(aShortFlagAction);
    assertFalse(bLongFlagAction.getRan());
    assertTrue(cShortParameterAction.getRan());
    CliTestUtilities.assertValidArguments(cShortParameterAction, "hello");
    assertFalse(dLongParameterAction.getRan());
  }

  @Test
  public void testEverything() throws IllegalArgumentException {
    cli.parse(new String[] { "-c", "hello", "--bob", "-a", "--dog", "world" });
    assertTrue(aShortFlagAction.getRan());
    CliTestUtilities.assertValidArguments(aShortFlagAction);
    assertTrue(bLongFlagAction.getRan());
    CliTestUtilities.assertValidArguments(bLongFlagAction);
    assertTrue(cShortParameterAction.getRan());
    CliTestUtilities.assertValidArguments(cShortParameterAction, "hello");
    assertTrue(dLongParameterAction.getRan());
    CliTestUtilities.assertValidArguments(dLongParameterAction, "world");
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
  public void usageTest() {
    System.out.println(cli.getUsage());
  }
}
