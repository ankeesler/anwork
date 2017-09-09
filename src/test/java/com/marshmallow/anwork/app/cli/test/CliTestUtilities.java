package com.marshmallow.anwork.app.cli.test;

import static org.junit.Assert.*;

/**
 * Some utilities to be used in the {@link CliTest}.
 *
 * @author Andrew
 * @date Sep 9, 2017
 */
public class CliTestUtilities {

  public static void assertValidArguments(TestCliAction action, String...expected) {
    String[] arguments = action.getArguments();
    assertNotNull(arguments.length);
    assertEquals(expected.length, arguments.length);
    for (int i = 0; i < arguments.length; i++) {
      assertEquals("The " + i + "th element of the arguments (" + arguments[i] + ")"
                   + " does not match the expected (" + expected[i] + ")",
                   arguments[i],
                   expected[i]);
    }
  }
}