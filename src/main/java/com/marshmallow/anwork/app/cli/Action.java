package com.marshmallow.anwork.app.cli;

/**
 * This is an action (i.e., some Java code) that can be run via a {@link Command}. When a
 * {@link Command} is issued on the command line, {@link #run(ArgumentValues, String[])} is called
 * on an object of this type.
 *
 * <p>
 * Created Sep 9, 2017
 * </p>
 *
 * @author Andrew
 */
public interface Action {

  /**
   * Run the CLI action.
   *
   * <p>The parameters to this method are two separate instances of {@link ArgumentValues}.
   * <ul>
   * <li>The first parameter contains any of the values that were passed with flags that have
   * {@link Argument}'s associated with them. The keys to this {@link ArgumentValues} data is the
   * short flag of the {@link Flag} in question (see {@link Flag#getShortFlag()}. Note that any
   * {@link Flag} without any {@link Argument} associated with it will exist in the
   * {@link ArgumentValues} data and will be paired with a {@link ArgumentType#BOOLEAN} type.</li>
   * <li>The second parameter contains any of the values that were passed to the {@link Command}
   * for which this {@link Action} is being run. The keys to this {@link ArgumentValues} data is
   * the name of the {@link Argument} (see {@link Argument#getName()}.</li>
   * </ul>
   *
   * <p>If there were three flags (-f (filename), -d, -n (number)) passed to this
   * {@link Command}, then a client implementation can access this flag data in the following way.
   * <pre>
   *   String filename = flags.getValue("f", ArgumentType.STRING);
   *   Boolean debug = flags.getValue("d", ArgumentType.DEBUG);
   *   Long number = flags.getValue("n", ArgumentType.NUMBER);
   * </pre>
   *
   * @param flags The {@link ArgumentValues} that represent the {@link Flag}'s {@link Argument}'s
   *     passed to this command
   * @param parameters The parameters to the action; this parameter is never <code>null</code>
   */
  public void run(ArgumentValues flags, String[] parameters);

}