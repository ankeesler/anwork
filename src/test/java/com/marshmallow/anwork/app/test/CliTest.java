package com.marshmallow.anwork.app.test;

import org.junit.Test;

import com.marshmallow.anwork.app.Cli;

import static org.junit.Assert.*;

import org.junit.Before;

/**
 * This is a test for the Cli.
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public class CliTest {

  private static final Cli.CliAction NULL_ACTION = (a) -> a.hashCode();

  private static class TestAction implements Cli.CliAction {

    private boolean ran = false;
    private String argument = null;

    @Override
    public void run(String argument) {
      ran = true;
      this.argument = argument;
    }

    public boolean getRan() {
      return ran;
    }

    public String getArgument() {
      return argument;
    }
  }

  private Cli cli = new Cli();
  private TestAction aShortFlagAction = new TestAction();
  private TestAction bLongFlagAction = new TestAction();
  private TestAction cShortArgumentAction = new TestAction();
  private TestAction dLongArgumentAction = new TestAction();

  @Before
  public void setupCli() {
    cli.addAction("a", null, "Description for flag a", null, aShortFlagAction);
    cli.addAction("b", "bob", null, null, bLongFlagAction);
    cli.addAction("c", null, "Do something related to something", "thing", cShortArgumentAction);
    cli.addAction("d", "dog", null, "name", dLongArgumentAction);
  }

  @Test(expected = IllegalArgumentException.class)
  public void addBadShortFlagTest() throws IllegalArgumentException {
    new Cli().addAction(null, "andrew", null, null, NULL_ACTION);
  }

  @Test(expected = IllegalArgumentException.class)
  public void addBadActionTest() throws IllegalArgumentException {
    new Cli().addAction("a", "andrew", null, null, null);
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadFlagSyntaxStart() throws IllegalArgumentException {
    cli.runActions(new String[] { "a", "-b" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testBadFlagSyntaxend() throws IllegalArgumentException {
    cli.runActions(new String[] { "-b", "a" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testUnknownShortFlag() throws IllegalArgumentException {
    cli.runActions(new String[] { "-b", "-z" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testUnknownLongFlag() throws IllegalArgumentException {
    cli.runActions(new String[] { "-b", "--zebra" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testNoShortArgument() throws IllegalArgumentException {
    cli.runActions(new String[] { "-b", "-c" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testNoLongArgument() throws IllegalArgumentException {
    cli.runActions(new String[] { "-b", "--dog" } );
  }

  @Test(expected = IllegalArgumentException.class)
  public void testLongFlagShortSyntax() throws IllegalArgumentException {
    cli.runActions(new String[] { "-b", "-bob" } );
  }

  @Test
  public void shortFlagOnlyTest() throws IllegalArgumentException {
    cli.runActions(new String[] { "-a" });
    assertTrue(aShortFlagAction.getRan());
    assertNull(aShortFlagAction.getArgument());
    assertFalse(bLongFlagAction.getRan());
    assertFalse(cShortArgumentAction.getRan());
    assertFalse(dLongArgumentAction.getRan());
  }

  @Test
  public void longFlagShortFlagTest() throws IllegalArgumentException {
    cli.runActions(new String[] { "--bob", "-a" });
    assertTrue(aShortFlagAction.getRan());
    assertNull(aShortFlagAction.getArgument());
    assertTrue(bLongFlagAction.getRan());
    assertNull(bLongFlagAction.getArgument());
    assertFalse(cShortArgumentAction.getRan());
    assertFalse(dLongArgumentAction.getRan());
  }

  @Test
  public void testShortArgumentShortFlag() throws IllegalArgumentException {
    cli.runActions(new String[] { "-c", "hello", "-a" });
    assertTrue(aShortFlagAction.getRan());
    assertNull(aShortFlagAction.getArgument());
    assertFalse(bLongFlagAction.getRan());
    assertTrue(cShortArgumentAction.getRan());
    assertEquals("hello", cShortArgumentAction.getArgument());
    assertFalse(dLongArgumentAction.getRan());
  }

  @Test
  public void testEverything() throws IllegalArgumentException {
    cli.runActions(new String[] { "-c", "hello", "--bob", "-a", "--dog", "world" });
    assertTrue(aShortFlagAction.getRan());
    assertNull(aShortFlagAction.getArgument());
    assertTrue(bLongFlagAction.getRan());
    assertNull(bLongFlagAction.getArgument());
    assertTrue(cShortArgumentAction.getRan());
    assertEquals("hello", cShortArgumentAction.getArgument());
    assertTrue(dLongArgumentAction.getRan());
    assertEquals("world", dLongArgumentAction.getArgument());
  }

  @Test
  public void testEmptyArgs() {
    cli.runActions(new String[] { });
    assertFalse(aShortFlagAction.getRan());
    assertFalse(bLongFlagAction.getRan());
    assertFalse(cShortArgumentAction.getRan());
    assertFalse(dLongArgumentAction.getRan());
  }

  @Test
  public void usageTest() {
    System.out.println(cli.getUsage());
  }
}
