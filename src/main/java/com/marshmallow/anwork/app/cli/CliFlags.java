package com.marshmallow.anwork.app.cli;

import java.util.HashMap;
import java.util.Map;

/**
 * This is a vague map-like class that gives {@link CliAction}'s an "easy" way to access
 * the flags that have been passed to a CLI command.
 *
 * <p>
 * Created Oct 3, 2017
 * </p>
 *
 * @author Andrew
 */
public class CliFlags {

  private final Map<String, Object> values = new HashMap<String, Object>();

  // Only allow this object to be created in the package.
  CliFlags() { }

  // This method is package-private and is meant to be used during parsing.
  void addShortFlagValue(String shortFlag, Object value) {
    values.put(shortFlag, value);
  }

  /**
   * Return all of the short flags that have values for them in this object.
   *
   * @return All of the short flags that have values for them in this object
   */
  public String[] getAllShortFlags() {
    return values.keySet().toArray(new String[0]);
  }

  /**
   * Get the value for a flag given the shortFlag (e.g., "d", "v", "o", etc.) and the expected type
   * of the value.
   *
   * @param shortFlag The shortFlag (e.g., "d", "v", "o", etc.)
   * @param valueType The expected {@link CliArgumentType} of the flag's value
   * @return The value for this short flag, or <code>null</code> if this flag was not passed to the
   *     command line
   * @throws IllegalArgumentException if this shortFlag is not associated with this
   *     expectedFlagType
   */
  public Object getValue(String shortFlag, CliArgumentType valueType)
      throws IllegalArgumentException {
    if (!values.containsKey(shortFlag)) {
      return null;
    }

    Object value = values.get(shortFlag);
    Class<?> conversionClass = valueType.getConverter().getConversionClass();
    if (!conversionClass.isInstance(value)) {
      throw new IllegalArgumentException("Incorrect flag value type. Got " + conversionClass
                                         + " but expected " + value.getClass());
    }
    return conversionClass.cast(value);
  }
}