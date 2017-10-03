package com.marshmallow.anwork.app.cli;

import java.util.IllegalFormatException;

/**
 * This is a static utility class that serves as an accessor for the instance
 * {@link CliArgumentConverter}'s.
 *
 * <p>
 * Note: if we want to allow multiple implementations of these converters in the future, this class
 * will need to go away, or become more flexible.
 * </p>
 *
 * <p>
 * Created Oct 3, 2017
 * </p>
 *
 * @author Andrew
 */
final class CliArgumentConverters {

  // This is a private utility class, so we hide the constructor.
  private CliArgumentConverters() { }

  private static CliArgumentConverter<Boolean> booleanConverter;
  private static CliArgumentConverter<String> stringConverter;
  private static CliArgumentConverter<Integer> integerConverter;

  /**
   * Get the instance {@link CliArgumentConverter} for {@link Boolean}'s.
   *
   * @return The instance {@link CliArgumentConverter} for {@link Boolean}'s
   */
  public static CliArgumentConverter<Boolean> getBooleanConverter() {
    if (booleanConverter == null) {
      booleanConverter = new CliArgumentConverter<Boolean>() {

        @Override
        public Class<Boolean> getConversionClass() {
          return Boolean.class;
        }

        @Override
        public Boolean convert(String string) throws IllegalFormatException {
          return Boolean.parseBoolean(string);
        }
      };
    }
    return booleanConverter;
  }

  /**
   * Get the instance {@link CliArgumentConverter} for {@link String}'s.
   *
   * @return The instance {@link CliArgumentConverter} for {@link String}'s
   */
  public static CliArgumentConverter<String> getStringConverter() {
    if (stringConverter == null) {
      stringConverter = new CliArgumentConverter<String>() {

        @Override
        public Class<String> getConversionClass() {
          return String.class;
        }

        @Override
        public String convert(String string) throws IllegalFormatException {
          return string;
        }
      };
    }
    return stringConverter;
  }

  /**
   * Get the instance {@link CliArgumentConverter} for {@link Integer}'s.
   *
   * @return The instance {@link CliArgumentConverter} for {@link Integer}'s
   */
  public static CliArgumentConverter<Integer> getIntegerConverter() {
    if (integerConverter == null) {
      integerConverter = new CliArgumentConverter<Integer>() {

        @Override
        public Class<Integer> getConversionClass() {
          return Integer.class;
        }

        @Override
        public Integer convert(String string) throws IllegalFormatException {
          return Integer.parseInt(string);
        }
      };
    }
    return integerConverter;
  }
}
