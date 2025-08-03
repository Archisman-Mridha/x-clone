import marimo

__generated_with = "0.14.10"
app = marimo.App(width="full", app_title="Post Labeller")


@app.cell
def _():
    post = """
    Spent an amazing night in the heart of Copenhagen, with 2 of my lovely colleagues ❤️
    @codingsafari and @johannesholzel.

    We talked primarily about CloudNative Distributed Systems, Data stream processing using Apache
    Flink and building Recommendation Engines.

    And dammnn! Those pastries were real good....😋

    #copenhagen #tech #distributed-systems #big-data #data-stream-processing
    """

    import re

    def preprocess_post(post):
        replacements = [
            # replace newline characters with simple spaces.
            (r'\r\n|\r|\n', ' '),

            # remove mentions and hashtags.
            (r'@[^\s]+', ''),
            (r'#[^\s]+', '')
        ]

        for replacement in replacements:
            old, new = replacement
            post = re.sub(old, new, post)

        return post

    preprocessed_post = preprocess_post(post)
    print(preprocessed_post)
    return (preprocessed_post,)


@app.cell
def _():
    import spacy

    # spaCy provides different models / trained pipelines, which can be installed as Python
    # packages.
    # We'll use the 'en_core_web_sm' model : English pipeline optimized for CPU.
    spacy_model = spacy.load("en_core_web_sm")
    return (spacy_model,)


@app.cell
def _(preprocessed_post, spacy_model):
    # Processing raw text intelligently is difficult: most words are rare, and it’s common for
    # words that look completely different to mean almost the same thing. The same words in a
    # different order can mean something completely different. Even splitting text into useful
    # word-like units can be difficult in many languages.
    #
    # While it’s possible to solve some problems starting from only the raw characters, it’s
    # usually better to use linguistic knowledge to add useful information. That’s exactly what
    # spaCy is designed to do:
    # you put in raw text, and get back a Doc object, that comes with a variety of annotations.
    doc = spacy_model(preprocessed_post)
    return (doc,)


@app.cell
def _(doc):
    # Linguistic annotations are available as Token attributes.
    #
    # NOTE : Like many NLP libraries, spaCy encodes all strings to hash values to reduce memory
    #        usage and improve efficiency. So to get the readable string representation of an
    #        attribute, we need to add an underscore _ to its name.
    for token in doc:
        print((token.text, token.lemma_, token.pos_, token.tag_, token.dep_, token.is_stop))
    return


@app.cell
def _(doc):
    def filter_tokens(doc):

        def filter_token(token):
            return (not token.is_alpha or
                    token.is_stop or
                    token.pos_ not in ["ADJ", "NOUN", "PROPN"] or
                    token.dep_ not in ["nsubj", "dobj", "pobj", "compound", "amod", "appos"])

        filtered_tokens = [token for token in doc if not filter_token(token)]
        return filtered_tokens

    filtered_tokens = filter_tokens(doc)

    for filtered_token in filtered_tokens:
        print((filtered_token.lemma_, filtered_token.pos_, filtered_token.tag_, filtered_token.dep_))
    return (filtered_tokens,)


@app.cell
def _(doc):
    # Noun chunks are ”base noun phrases” – flat phrases that have a noun as their head. You can
    # think of noun chunks as a noun plus the words describing the noun – for example :
    # “the lavish green grass” or “the world’s largest tech fund”. 
    for noun_chunk in doc.noun_chunks:
        print((noun_chunk.text, noun_chunk.root.text, noun_chunk.root.dep_))
    return


@app.cell
def _(doc):
    # spaCy features an extremely fast statistical entity recognition system, that assigns labels
    # to contiguous spans of tokens. The default trained pipelines can identify a variety of named
    # and numeric entities, including companies, locations, organizations and products.
    for entity in doc.ents:
        print((entity.text, entity.label_))
    return


@app.cell
def _(doc):
    def filter_entities(doc):

        def filter_entity(entity):
            return (entity.label_ not in ["ORG", "PRODUCT"])

        filtered_entities = [entity for entity in doc.ents if not filter_entity(entity)]
        return filtered_entities

    filtered_entities = filter_entities(doc)

    for filtered_entity in filtered_entities:
        print((filtered_entity.text, filtered_entity.label_))
    return (filtered_entities,)


@app.cell
def _(filtered_entities, filtered_tokens):
    post_labels = list(set(
        [filtered_token.text.lower() for filtered_token in filtered_tokens] +
        [filtered_entity.text.lower() for filtered_entity in filtered_entities]))

    print(post_labels)
    return (post_labels,)


@app.cell
def _(post_labels, preprocessed_post):
    import keybert
    from sklearn.feature_extraction.text import CountVectorizer
    from sentence_transformers import SentenceTransformer

    # Sentence Transformers (a.k.a. SBERT) is the go-to Python module for accessing, using, and
    # training state-of-the-art embedding and reranker models. It can be used to compute embeddings
    # using Sentence Transformer models, to calculate similarity scores using Cross-Encoder
    # (a.k.a. reranker) models, or to generate sparse embeddings using Sparse Encoder models. This
    # unlocks a wide range of applications, including semantic search, semantic textual similarity,
    # and paraphrase mining.
    #
    # A wide selection of over 10,000 pre-trained Sentence Transformers models are available for
    # immediate use on 🤗 Hugging Face.
    sentence_transformer_model = SentenceTransformer("all-MiniLM-L6-v2")

    # KeyBERT supports quite a few embedding models. Having the option to choose embedding models
    # allow you to leverage pre-trained embeddings that suit your use-case.
    keyBERT_model = keybert.KeyBERT(model=sentence_transformer_model)

    def rank_post_labels(post_labels: str) -> str:
        # An unexpectly important component of KeyBERT is the CountVectorizer. In KeyBERT, it is
        # used to split up your documents into candidate keywords and keyphrases.
        # Since we use the vectorizer to split up the documents after embedding them, we can parse
        # the document however we want as it does not affect the quality of the document embeddings.
        vectorizer = CountVectorizer(ngram_range=(1, 3),
                                     vocabulary=post_labels)

        # As a default, KeyBERT simply compares the documents and candidate keywords/keyphrases
        # based on their cosine similarity. However, this might lead to very similar words ending
        # up in the list of most accurate keywords/keyphrases. To make sure they are a bit more
        # diversified we can use Maximal Margin Relevance (MMR).
        keywords = keyBERT_model.extract_keywords(preprocessed_post.lower(),
                                                  vectorizer=vectorizer,
                                                  top_n=5,
                                                  use_mmr=True,
                                                  diversity=0.5)

        ranked_post_labels = [keyword[0] for keyword in keywords]
        return ranked_post_labels

    ranked_post_labels = rank_post_labels(post_labels)
    print(ranked_post_labels)
    return ranked_post_labels, sentence_transformer_model


@app.cell
def _(ranked_post_labels, sentence_transformer_model):
    post_categories = [
        "travel",
        "technology",
        "photography",
        "music",
        "food",
        "football",
        "dance",
        "cricket",
        "yoga",
        "gym",
        "art"
    ]

    from sentence_transformers import util
    from collections import Counter

    def match_ranked_post_labels_to_post_categories(ranked_post_labels):
        ranked_post_label_encodings = sentence_transformer_model.encode(ranked_post_labels,
                                                                           convert_to_tensor=True)
        post_category_encodings = sentence_transformer_model.encode(post_categories, convert_to_tensor=True)

        cosine_similarities = util.cos_sim(ranked_post_label_encodings, post_category_encodings)

        post_category_match_counts = Counter()
        for i, ranked_post_label in enumerate(ranked_post_labels):
            cosine_similarity = cosine_similarities[i]

            best_matched_post_category_index = cosine_similarity.argmax().item()
            best_matched_post_category = post_categories[best_matched_post_category_index]

            print({
                "post label": ranked_post_label,
                "matched category" : best_matched_post_category,
                "similarity score" : cosine_similarity.max().item()
            })

            post_category_match_counts[best_matched_post_category] += 1

        total_match_count = sum(post_category_match_counts.values())
        post_category_match_percentages = {
            post_category:
                round(post_category_match_count * 100 / total_match_count, 2)

            for post_category, post_category_match_count in post_category_match_counts.items()
        }
        return post_category_match_percentages

    match_percentages = match_ranked_post_labels_to_post_categories(ranked_post_labels)
    print(match_percentages)
    return


if __name__ == "__main__":
    app.run()
