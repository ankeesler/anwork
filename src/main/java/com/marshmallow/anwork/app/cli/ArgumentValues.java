package com.marshmallow.anwork.app.cli;

import java.util.HashMap;
import java.util.Map;

/**
 * This is a vague map-like class that gives {@link Action}'s an "easy" way to access
 * the {@link Flag}'s/{@link Command}'s {@link Argument} values that have been passed on the
 * command line.
 *
 * <p>
 * Created Oct 3, 2017
 * </p>
 *
 * @author Andrew
 */
public class ArgumentValues {

  private final Map<String, Object> values = new HashMap<String, Object>();

  // Only allow this object to be created in the package.
  ArgumentValues() { }

  // This method is package-private and is meant to be used during parsing.
  void addShortFlagValue(String shortFlag, Object value) {
    values.put(shortFlag, value);
  }

  /**
   * Return all of the keys associated with the values in this class.
   *
   * @return All of the keys associated with the values in this class
   */
  public String[] getAllKeys() {
    return values.keySet().toArray(new String[0]);
  }

  /**
   * Return whether or not there is a value in this class that is associated with the provided key.
   *
   * @param key The key that may or may not be associated with a value in this class
   * @return whether or not there is a value in this class that is associated with the provided key
   */
  public boolean containsKey(String key) {
    return values.containsKey(key);
  }

  /**
   * Get the value for a argument given the key and the expected type of the value.
   *
   * @param <T> the type of the backing Java type for this value's {@link ArgumentType}
   * @param key The key for the value; this key is context dependent
   * @param valueType The expected {@link ArgumentType} of the flag's value
   * @return The value for this key, or <code>null</code> if no value for this key was passed to
   *     the command line
   * @throws IllegalArgumentException if this key is not associated with this valueType
   */
  public <T> T getValue(String key, ArgumentType<T> valueType)
      throws IllegalArgumentException {
    if (!values.containsKey(key)) {
      return null;
    }

    Object value = values.get(key);
    Class<T> conversionClass = valueType.getConversionClass();
    if (!conversionClass.isInstance(value)) {
      throw new IllegalArgumentException("Incorrect flag value type. Got " + conversionClass
                                         + " but expected " + value.getClass());
    }
    return conversionClass.cast(value);
  }
}