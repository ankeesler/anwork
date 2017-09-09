package com.marshmallow.anwork.app.test;

import org.junit.Test;

import com.marshmallow.anwork.app.CliNode;
import com.marshmallow.anwork.app.CliAction;

import static org.junit.Assert.*;

import org.junit.Before;

/**
 * This is a test for the Cli.
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public class CliTest {

  private static final CliAction NULL_ACTION = (a) -> a.hashCode();

  private static class TestAction implements CliAction {

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

  private CliNode cli = new CliNode("whatever");
  private TestAction aShortFlagAction = new TestAction();
  private TestAction bLongFlagAction = new TestAction();
  private TestAction cShortArgumentAction = new TestAction();
  private TestAction dLongArgumentAction = new TestAction();

  @Before
  public void setupCli() {
    cli.addFlag("a", null, "Description for flag a", null, aShortFlagAction);
    cli.addFlag("b", "bob", null, null, bLongFlagAction);
    cli.addFlag("c", null, "Do something related to something", "thing", cShortArgumentAction);
    cli.addFlag("d", "dog", null, "name", dLongArgumentAction);
  }

  @Test(expected = IllegalArgumentException.class)
  public void addBadShortFlagTest() throws IllegalArgumentException {
    new CliNode("whatever").addFlag(null, "andrew", null, null, NULL_ACTION);
  }

  @Test(expected = IllegalArgumentException.class)
  public void addBadActionTest() throws IllegalArgumentException {
    new CliNode("whatever").addFlag("a", "andrew", null, null, null);
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
