package com.marshmallow.anwork.app.cli;

/**
 * This is a form of documentation that can be used to document a {@link Cli} API.
 *
 * <p>
 * Created Oct 16, 2017
 * </p>
 *
 * @author Andrew
 */
public enum DocumentationType {
  TEXT(TextDocumentationGenerator.class),
  GITHUB_MARKDOWN(GithubReadmeDocumentationGenerator.class),
  ;

  private final Class<?extends DocumentationGenerator> clazz;

  private DocumentationType(Class<? extends DocumentationGenerator> clazz) {
    this.clazz = clazz;
  }

  /**
   * Create an instance of the {@link Class} associated with this {@link DocumentationType}.
   *
   * @return An instance of the {@link Class} associated with this {@link DocumentationType}
   * @throws Exception if the instantiation of the {@link Class} for this {@link DocumentationType}
   *     fails
   */
  DocumentationGenerator instantiate() throws Exception {
    return clazz.newInstance();
  }
}
