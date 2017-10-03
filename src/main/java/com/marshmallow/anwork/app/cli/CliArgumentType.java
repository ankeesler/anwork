package com.marshmallow.anwork.app.cli;

/**
 * This is a list of the allowed types in this CLI system.
 *
 * <p>
 * Each type is paired with a {@link CliArgumentConverter} which determines the Java class to which
 * the {@link CliArgumentType} maps.
 * </p>
 *
 * <p>
 * Created Oct 3, 2017
 * </p>
 *
 * @author Andrew
 */
public enum CliArgumentType {
  BOOLEAN(CliArgumentConverters.getBooleanConverter()),
  STRING(CliArgumentConverters.getStringConverter()),
  INTEGER(CliArgumentConverters.getIntegerConverter()),
  ;

  private final CliArgumentConverter<?> converter;

  private CliArgumentType(CliArgumentConverter<?> converter) {
    this.converter = converter;
  }

  /**
   * Get the {@link CliArgumentConverter} associated with this {@link CliArgumentType}.
   *
   * @return The {@link CliArgumentConverter} associated with this {@link CliArgumentType}
   */
  public CliArgumentConverter<?> getConverter() {
    return converter;
  }
}