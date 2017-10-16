package com.marshmallow.anwork.app.cli;

import java.io.PrintWriter;

/**
 * This is an object that generates documentation for a {@link Cli} API in some format (see
 * {@link DocumentationType} for formats).
 *
 * <p>
 * Created Oct 16, 2017
 * </p>
 *
 * @author Andrew
 */
public interface DocumentationGenerator {

  /**
   * Get the {@link DocumentationType} that this {@link DocumentationGenerator} produces.
   *
   * @return The {@link DocumentationType} that this {@link DocumentationGenerator} produces
   */
  public DocumentationType getType();

  /**
   * Generate the documentation for a {@link Cli} API.
   *
   * @param cli The {@link Cli} for which documentation will be generated
   * @param writer The {@link PrintWriter} that the documentation text will be written to
   */
  public void generate(Cli cli, PrintWriter writer);
}
