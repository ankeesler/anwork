package com.marshmallow.anwork.app.cli.test;

import static org.junit.Assert.assertNotNull;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.DocumentationGenerator;
import com.marshmallow.anwork.app.cli.DocumentationGeneratorFactory;
import com.marshmallow.anwork.app.cli.DocumentationType;
import com.marshmallow.anwork.app.cli.Visitor;

import java.io.IOException;
import java.io.PrintWriter;
import java.io.Writer;

import org.junit.Before;
import org.junit.Test;

/**
 * This is a generic abstract base class that creates a {@link Cli} and runs some tests.
 *
 * <p>
 * Created Oct 4, 2017
 * </p>
 *
 * @author Andrew
 */
public abstract class BaseCliTest {

  /**
   * This class is a {@link Writer} that does nothing. It is useful for writing controlled tests.
   *
   * <p>
   * Created Oct 16, 2017
   * </p>
   *
   * @author Andrew
   */
  private static class NoOpWriter extends Writer {
    @Override
    public void write(char[] cbuf, int off, int len) throws IOException {
      // no-op
    }

    @Override
    public void flush() throws IOException {
      // no-op
    }

    @Override
    public void close() throws IOException {
      // no-op
    }
  }

  private Cli cli;

  @Before
  public void setupCli() throws Exception {
    cli = createCli();
    assertNotNull(cli);
  }

  /**
   * Create the {@link Cli} object that will be used in tests.
   *
   * @return The {@link Cli} object that will be used in tests
   * @throws Exception if the {@link Cli} could not be created
   */
  protected abstract Cli createCli() throws Exception;

  /**
   * Run the provided arguments through the {@link Cli} created through {@link #createCli()}.
   *
   * @param arguments The arguments to run through the {@link Cli} via {@link Cli#parse(String[])}
   */
  protected void parse(String...arguments) {
    cli.parse(arguments);
  }

  /**
   * Run a {@link Visitor} on the {@link Cli} created through {@link #createCli()}.
   *
   * @param visitor The {@link Visitor} to visit the {@link Cli} tree created through
   *     {@link #createCli()}
   */
  protected void visit(Visitor visitor) {
    cli.visit(visitor);
  }

  /**
   * This test validates that nothing bad happens when {@link Cli#getUsage()} is run...
   */
  @Test
  public void usageTest() {
    System.out.println(cli.getUsage());
  }

  /**
   * This test ensures that the {@link Cli} instance for this class can successfully have
   * documentation generated for it via a {@link DocumentationGenerator}.
   */
  @Test
  public void documentationGenerationTest() throws Exception {
    for (DocumentationType documentationType : DocumentationType.values()) {
      DocumentationGenerator generator
          = DocumentationGeneratorFactory.getInstance().createGenerator(documentationType);
      try {
        generator.generate(cli, new PrintWriter(new NoOpWriter(), false));
      } catch (Exception e) {
        throw new Exception("Documentation generation failed for type " + documentationType);
      }
    }
  }
}
