package com.marshmallow.anwork.app.cli.test;

import static org.junit.Assert.assertNotNull;

import com.marshmallow.anwork.app.cli.Cli;
import com.marshmallow.anwork.app.cli.CliVisitor;

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
   * Run a {@link CliVisitor} on the {@link Cli} created through {@link #createCli()}.
   *
   * @param visitor The {@link CliVisitor} to visit the {@link Cli} tree created through
   *     {@link #createCli()}
   */
  protected void visit(CliVisitor visitor) {
    cli.visit(visitor);
  }

  /**
   * This test validates that nothing bad happens when {@link Cli#getUsage()} is run...
   */
  @Test
  public void usageTest() {
    System.out.println(cli.getUsage());
  }
}
