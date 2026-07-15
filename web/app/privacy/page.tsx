import type { Metadata } from "next";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";

// ---------------------------------------------------------------------------
// PLACEHOLDERS — replace every value below before publishing this page.
// They are referenced throughout the policy so a single edit updates all
// occurrences.
// ---------------------------------------------------------------------------
const COMPANY_NAME = "[Company Name]";
const WEBSITE_URL = "[Website URL]";
const CONTACT_EMAIL = "[Contact Email]";
const BUSINESS_ADDRESS = "[Business Address]";
const EFFECTIVE_DATE = "[Effective Date]";

export const metadata: Metadata = {
  title: "Privacy Policy — ThumbnailIQ",
  description:
    "Learn how ThumbnailIQ collects, uses, shares, and protects your personal information, and the choices and rights you have over your data.",
  alternates: { canonical: "/privacy" },
};

// The policy is data-driven: each section is a heading plus JSX body. This
// keeps the TOC, anchor ids, and heading hierarchy generated from one list so
// they can never drift out of sync.
interface PolicySection {
  id: string;
  title: string;
  body: React.ReactNode;
}

// Consistent styling for inline lists inside sections.
function Bullets({ items }: { items: React.ReactNode[] }) {
  return (
    <ul className="list-disc space-y-2 pl-6">
      {items.map((item, i) => (
        <li key={i}>{item}</li>
      ))}
    </ul>
  );
}

const SECTIONS: PolicySection[] = [
  {
    id: "introduction",
    title: "Introduction",
    body: (
      <>
        <p>
          {COMPANY_NAME} (&quot;we&quot;, &quot;us&quot;, or &quot;our&quot;) operates the ThumbnailIQ
          service available at {WEBSITE_URL} (the &quot;Service&quot;). This Privacy Policy explains
          what information we collect, how we use and share it, and the choices and rights you have
          over it.
        </p>
        <p>
          By creating an account or otherwise using the Service, you agree to the collection and use
          of information as described in this policy. If you do not agree, please do not use the
          Service.
        </p>
      </>
    ),
  },
  {
    id: "information-we-collect",
    title: "Information We Collect",
    body: (
      <>
        <h3 className="text-base font-semibold text-white">Information you provide to us</h3>
        <Bullets
          items={[
            <>
              <strong className="text-gray-200">Account information</strong> — your name, email
              address, and password when you register.
            </>,
            <>
              <strong className="text-gray-200">Content you upload</strong> — thumbnail images,
              keywords, video titles, and related material you submit for analysis.
            </>,
            <>
              <strong className="text-gray-200">Workspace and team data</strong> — workspace names,
              member invitations, and roles you configure.
            </>,
            <>
              <strong className="text-gray-200">Payment information</strong> — when you purchase a
              subscription, payments are processed by our payment providers; we receive confirmation
              of payment and limited billing details, but we never store full card numbers.
            </>,
            <>
              <strong className="text-gray-200">Communications</strong> — messages you send us, such
              as support requests.
            </>,
          ]}
        />
        <h3 className="text-base font-semibold text-white">Information collected automatically</h3>
        <Bullets
          items={[
            <>
              <strong className="text-gray-200">Usage data</strong> — pages visited, features used,
              analyses run, and interactions with the Service.
            </>,
            <>
              <strong className="text-gray-200">Device and log data</strong> — IP address, browser
              type and version, operating system, timestamps, and error logs. We also record the IP
              address and device details of sign-ins to help keep your account secure.
            </>,
          ]}
        />
      </>
    ),
  },
  {
    id: "how-we-use-information",
    title: "How We Use Information",
    body: (
      <>
        <p>We use the information we collect to:</p>
        <Bullets
          items={[
            "Provide, operate, and maintain the Service, including analyzing the thumbnails you submit and generating scores and recommendations.",
            "Create and manage your account, workspaces, and team memberships.",
            "Process subscription payments and send transactional emails such as verification codes, password resets, sign-in alerts, and billing confirmations.",
            "Monitor usage against your plan limits and prevent abuse, fraud, and unauthorized access.",
            "Improve the Service, develop new features, and fix bugs.",
            "Comply with legal obligations and enforce our terms.",
          ]}
        />
      </>
    ),
  },
  {
    id: "cookies",
    title: "Cookies and Tracking Technologies",
    body: (
      <>
        <p>
          We use cookies and similar browser storage technologies (such as local storage) to keep
          you signed in, remember preferences like your selected workspace, and understand how the
          Service is used.
        </p>
        <p>
          You can control cookies through your browser settings, including blocking or deleting
          them. Disabling essential cookies or storage may prevent parts of the Service — such as
          staying signed in — from working correctly.
        </p>
      </>
    ),
  },
  {
    id: "analytics",
    title: "Analytics",
    body: (
      <p>
        We may use analytics tools to collect aggregated, non-identifying information about how
        visitors use the Service — for example, which features are most used and how users move
        through the app. This helps us prioritize improvements. Where analytics providers are used,
        their collection is governed by their own privacy policies, and we configure them to
        minimize the personal data they receive.
      </p>
    ),
  },
  {
    id: "third-party-services",
    title: "Third-Party Services",
    body: (
      <>
        <p>
          We rely on a small number of service providers to operate the Service. These providers
          process data only on our behalf and under contractual safeguards:
        </p>
        <Bullets
          items={[
            <>
              <strong className="text-gray-200">Payment processors</strong> (e.g. Razorpay, Stripe)
              — to handle subscription payments securely.
            </>,
            <>
              <strong className="text-gray-200">Cloud hosting and storage</strong> (e.g. Amazon Web
              Services) — to host the Service and store uploaded images.
            </>,
            <>
              <strong className="text-gray-200">Email delivery</strong> — to send transactional
              emails such as verification codes and account notifications.
            </>,
            <>
              <strong className="text-gray-200">AI analysis providers</strong> — thumbnail images
              and related metadata you submit may be processed by machine-learning services to
              generate scores and suggestions.
            </>,
            <>
              <strong className="text-gray-200">Public platform data</strong> (e.g. the YouTube Data
              API) — to retrieve publicly available competitor video information for comparisons.
            </>,
          ]}
        />
      </>
    ),
  },
  {
    id: "data-sharing",
    title: "Data Sharing",
    body: (
      <>
        <p>We do not sell your personal information. We share information only:</p>
        <Bullets
          items={[
            "With the service providers listed above, to the extent needed to operate the Service.",
            "Within your workspace — your name and email address are visible to other members of workspaces you join.",
            "When required by law, regulation, legal process, or an enforceable governmental request.",
            "To protect the rights, property, or safety of our users, the public, or the Service.",
            "In connection with a merger, acquisition, or sale of assets, in which case we will notify you before your information becomes subject to a different privacy policy.",
          ]}
        />
      </>
    ),
  },
  {
    id: "data-security",
    title: "Data Security",
    body: (
      <p>
        We take reasonable technical and organizational measures to protect your information,
        including encrypted connections (HTTPS/TLS), hashed passwords, hashed verification codes,
        scoped access controls between workspaces, and sign-in alerts for your account. No method of
        transmission or storage is completely secure, however, and we cannot guarantee absolute
        security. If we become aware of a breach affecting your personal data, we will notify you as
        required by applicable law.
      </p>
    ),
  },
  {
    id: "data-retention",
    title: "Data Retention",
    body: (
      <p>
        We retain your information for as long as your account is active or as needed to provide the
        Service. Uploaded thumbnails and analysis results are kept so you can revisit and compare
        them. If you delete your account, we delete or anonymize your personal information within a
        reasonable period, except where we must retain it to comply with legal, tax, or accounting
        obligations, resolve disputes, or enforce agreements.
      </p>
    ),
  },
  {
    id: "your-rights",
    title: "Your Rights",
    body: (
      <>
        <p>
          Depending on where you live, you may have some or all of the following rights regarding
          your personal information:
        </p>
        <Bullets
          items={[
            "Access — request a copy of the personal information we hold about you.",
            "Correction — ask us to correct inaccurate or incomplete information.",
            "Deletion — ask us to delete your personal information.",
            "Portability — receive your information in a structured, machine-readable format.",
            "Objection and restriction — object to or ask us to limit certain processing.",
            "Withdraw consent — where processing is based on consent, withdraw it at any time.",
          ]}
        />
        <p>
          To exercise any of these rights, contact us at {CONTACT_EMAIL}. We will respond within the
          timeframe required by applicable law. You may also have the right to lodge a complaint
          with your local data protection authority.
        </p>
      </>
    ),
  },
  {
    id: "childrens-privacy",
    title: "Children's Privacy",
    body: (
      <p>
        The Service is not directed to children under the age of 13 (or the equivalent minimum age
        in your jurisdiction), and we do not knowingly collect personal information from children.
        If you believe a child has provided us with personal information, please contact us at{" "}
        {CONTACT_EMAIL} and we will delete it promptly.
      </p>
    ),
  },
  {
    id: "international-transfers",
    title: "International Data Transfers",
    body: (
      <p>
        Our service providers may store and process information in countries other than your own,
        which may have different data protection laws. Where we transfer personal data
        internationally, we rely on appropriate safeguards such as standard contractual clauses or
        equivalent mechanisms recognized by applicable law, and we require providers to protect your
        information consistently with this policy.
      </p>
    ),
  },
  {
    id: "changes",
    title: "Changes to This Privacy Policy",
    body: (
      <p>
        We may update this Privacy Policy from time to time. When we make material changes, we will
        update the effective date at the top of this page and, where appropriate, notify you by
        email or through the Service. Your continued use of the Service after changes take effect
        constitutes acceptance of the revised policy.
      </p>
    ),
  },
  {
    id: "contact",
    title: "Contact Information",
    body: (
      <>
        <p>If you have questions or concerns about this Privacy Policy or our data practices, contact us:</p>
        <address className="not-italic">
          <p className="font-medium text-gray-200">{COMPANY_NAME}</p>
          <p>{BUSINESS_ADDRESS}</p>
          <p>
            Email:{" "}
            <a
              href={`mailto:${CONTACT_EMAIL}`}
              className="rounded text-brand-300 transition-colors hover:text-brand-200 focus-visible:outline focus-visible:outline-2 focus-visible:outline-brand-400"
            >
              {CONTACT_EMAIL}
            </a>
          </p>
        </address>
      </>
    ),
  },
];

export default function PrivacyPolicyPage() {
  return (
    <div className="flex min-h-screen flex-col">
      <Navbar />

      <main className="mx-auto w-full max-w-3xl flex-1 px-6 py-16">
        <header className="mb-10">
          <h1 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">
            Privacy Policy
          </h1>
          <p className="mt-3 text-sm text-gray-500">Effective date: {EFFECTIVE_DATE}</p>
        </header>

        {/* Table of contents: in-page anchors, fully keyboard-navigable. */}
        <nav
          aria-label="Table of contents"
          className="mb-12 rounded-xl border border-surface-300 bg-surface-100 p-5"
        >
          <h2 className="mb-3 text-sm font-semibold uppercase tracking-widest text-gray-400">
            On this page
          </h2>
          <ol className="grid gap-x-8 gap-y-1.5 text-sm sm:grid-cols-2">
            {SECTIONS.map((s, i) => (
              <li key={s.id}>
                <a
                  href={`#${s.id}`}
                  className="inline-block rounded py-0.5 text-gray-400 transition-colors hover:text-brand-300 focus-visible:outline focus-visible:outline-2 focus-visible:outline-brand-400"
                >
                  {i + 1}. {s.title}
                </a>
              </li>
            ))}
          </ol>
        </nav>

        <div className="space-y-12">
          {SECTIONS.map((s, i) => (
            <section key={s.id} id={s.id} aria-labelledby={`${s.id}-heading`} className="scroll-mt-24">
              <h2
                id={`${s.id}-heading`}
                className="mb-4 text-xl font-semibold tracking-tight text-white"
              >
                {i + 1}. {s.title}
              </h2>
              <div className="space-y-4 text-sm leading-relaxed text-gray-300">{s.body}</div>
            </section>
          ))}
        </div>
      </main>

      <Footer />
    </div>
  );
}
