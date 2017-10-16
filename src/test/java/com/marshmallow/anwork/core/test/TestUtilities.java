package com.marshmallow.anwork.core.test;

import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.fail;

import java.io.File;
import java.net.URL;

/**
 * A class that holds test utilties.
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public final class TestUtilities {

  /**
   * Get a {@link File} with the passed name in the same package as the passed {@link Class}.
   *
   * <p>
   * Note: this method will call {@link Assert#fail(String)} if the file cannot be found.
   * </p>
   *
   * @param name The name of the file
   * @param clazz The {@link Class} to use to get the file
   * @return A {@link File} with the passed name in the same package as the passed {@link Class}
   */
  public static File getFile(String name, Class<?> clazz) {
    URL url = clazz.getResource(name);
    try {
      File file = new File(url.toURI());
      System.out.println("Loaded file with name " + name + " and class " + clazz + ": " + file);
      return file;
    } catch (Exception e) {
      fail("Cannot convert " + name + " to file for class " + clazz + ": " + e);
      return null;
    }
  }

  /**
   * Assert that an array is equal to another array passed as a variadic argument.
   *
   * @param array The array under test
   * @param expected The expected array
   */
  public static void assertVariadicArrayEquals(Object[] array, Object...expected) {
    assertArrayEquals(expected, array);
  }
}
