package com.marshmallow.anwork.app.cli;

/**
 * This represents the type of an {@link Argument}. Is it generic to a Java type that can be used
 * to access this type at runtime.
 *
 * <p>
 * For example, a command line {@link Argument} may represent a number. The {@link ArgumentType}
 * {@link ArgumentType#NUMBER} denotes that this argument in question is a number.
 * </p>
 *
 * <p>
 * Each {@link ArgumentType} implements the {@link ArgumentType#convert} method which takes a
 * {@link String} command line argument and turns it into an instance of the backing Java type for
 * this class.
 * </p>
 *
 * <p>
 * Created Oct 14, 2017
 * </p>
 *
 * @author Andrew
 */
public interface ArgumentType<T> {

  /**
   * This is a default implementation of a {@link ArgumentType} for a string {@link Argument}.
   */
  public static final ArgumentType<String> STRING = new ArgumentType<String>() {
    @Override
    public String convert(String string) {
      return string;
    }

    @Override
    public Class<String> getConversionClass() {
      return String.class;
    }
  };

  /**
   * This is a default implementation of a {@link ArgumentType} for a number {@link Argument}.
   */
  public static final ArgumentType<Long> NUMBER = new ArgumentType<Long>() {
    @Override
    public Long convert(String string) {
      return Long.parseLong(string);
    }

    @Override
    public Class<Long> getConversionClass() {
      return Long.class;
    }
  };

  /**
   * This is a default implementation of a {@link ArgumentType} for a boolean {@link Argument}.
   */
  public static final ArgumentType<Boolean> BOOLEAN = new ArgumentType<Boolean>() {
    @Override
    public Boolean convert(String string) {
      return Boolean.parseBoolean(string);
    }

    @Override
    public Class<Boolean> getConversionClass() {
      return Boolean.class;
    }
  };

  /**
   * Turn a {@link String} from the command line into an instance of the backing Java type for this
   * class.
   *
   * @param string The {@link String} from the command line that represents an argument of this
   *     type
   * @return An instance of the backing Java type for this class represented by the provided
   *     {@link String}
   * @throws IllegalArgumentException if the string is not able to be turned into an instance of
   *     the backing Java type
   */
  public T convert(String string) throws IllegalArgumentException;

  /**
   * This is a utility method that returns the {@link Class} that represents the Java runtime type
   * for this {@link ArgumentType}.
   *
   * @return The {@link Class} that represents the Java runtime type for this {@link ArgumentType}
   */
  public Class<T> getConversionClass();
}
